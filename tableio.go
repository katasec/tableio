package tableio

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/gertd/go-pluralize"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"
	"github.com/katasec/tableio/reflectx"
	"github.com/katasec/utils/errx"
	_ "github.com/mattn/go-sqlite3"
)

type TableIO[T any] struct {
	DB        *sql.DB
	tableName string

	// This is a comma separated list of fields for the table used for select statements
	selectList string

	// dbFields is a list of fields in the struct that have a 'db' tag
	dbFields []reflectx.FieldInfo
}

func NewTableIO[T any](driverName string, dataSourceName string) (*TableIO[T], error) {
	//db, err := sqlx.Connect(driverName, dataSourceName)
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	tableio := &TableIO[T]{
		DB:        db,
		tableName: GetTableName[T](),
		dbFields:  reflectx.GetDbStructFields[T](),
	}

	// Initialize the 'select list' for table
	tableio.selectList = tableio.genSelectList()

	return tableio, nil
}

func (me *TableIO[T]) Insert(data T) error {

	sqlCmd := "insert into " + me.tableName + "(" + me.selectList + ") values (" + reflectx.GetStructValues(data) + ")"

	// Run Insert
	_, err := me.DB.Exec(sqlCmd)
	errx.PanicOnError(err)

	// return data
	return nil
}

func (me *TableIO[T]) InsertMany(data []T) error {

	// Gen Sql Command
	sqlCmd := ""
	for _, item := range data {
		sqlCmd += "insert into " + me.tableName + "(" + me.selectList + ") values (" + reflectx.GetStructValues(item) + "); \n"
	}

	// Run Insert
	_, err := me.DB.Exec(sqlCmd)
	errx.PanicOnError(err)

	// return data
	return nil
}

func (me *TableIO[T]) CreateTableIfNotExists(verbose ...bool) error {
	var debug bool

	if len(verbose) > 0 {
		debug = verbose[0]
	}

	var sb strings.Builder

	tableName := GetTableName[T]()

	// Start Create Table Command
	sb.WriteString("CREATE TABLE IF NOT EXISTS " + tableName + " (\n")

	// Add fields
	sb.WriteString(reflectx.GenSqlForFields(me.dbFields))

	// End Command
	sb.WriteString(");")

	// Generate string
	sqlCmd := sb.String()
	if debug {
		fmt.Println(sqlCmd)
	}

	//Execute SQL to create table
	_, err := me.DB.Exec(sqlCmd)
	errx.PanicOnErrorf(err, "Error creating table %s", tableName)
	if err == nil {
		fmt.Println("Create table '" + me.tableName + "' successfully")
	}
	return nil
}

func (me *TableIO[T]) DeleteTableIfExists() {

	tableName := GetTableName[T]()

	// Start Create Table Commands
	sqlCmd := "DROP TABLE IF EXISTS " + tableName + ";"

	// Execute SQL to create table
	_, err := me.DB.Exec(sqlCmd)
	errx.PanicOnError(err)
	fmt.Println("Deleted table '" + me.tableName + "' successfully.")

}

func (me *TableIO[T]) Close() {
	me.DB.Close()
}

func (me *TableIO[T]) All() []T {
	//var data T

	// Construct select statement
	sqlCmd := "select " + me.selectList + " from " + me.tableName
	//fmt.Println(sqlCmd)

	// Run Query
	rows, err := me.DB.Query(sqlCmd)
	errx.PanicOnError(err)

	// Get column types and count
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		errx.PanicOnError(err)
	}
	count := len(columnTypes)
	finalRows := []interface{}{}

	// Loop through rows
	for rows.Next() {

		// Create pointers to each column of appropriate type
		scanArgs := make([]interface{}, count)
		for i, v := range columnTypes {
			switch v.DatabaseTypeName() {
			case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
				scanArgs[i] = new(sql.NullString)
			case "BOOL":
				scanArgs[i] = new(sql.NullBool)
			case "INT4":
				scanArgs[i] = new(sql.NullInt64)
			default:
				scanArgs[i] = new(sql.NullString)
			}
		}

		// Scan row into pointers
		err := rows.Scan(scanArgs...)
		errx.PanicOnError(err)

		masterData := map[string]interface{}{}
		for i, v := range columnTypes {

			if z, ok := (scanArgs[i]).(*sql.NullBool); ok {
				masterData[v.Name()] = z.Bool
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullString); ok {
				masterData[v.Name()] = z.String
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt64); ok {
				masterData[v.Name()] = z.Int64
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullFloat64); ok {
				masterData[v.Name()] = z.Float64
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt32); ok {
				masterData[v.Name()] = z.Int32
				continue
			}

			masterData[v.Name()] = scanArgs[i]
		}

		finalRows = append(finalRows, masterData)
	}

	z, err := json.Marshal(finalRows)
	errx.PanicOnError(err)
	//fmt.Println(string(z))

	var data []T
	json.Unmarshal(z, &data)
	return data
}

// genSelectList returns a comma separated list of fields for the table used for select statements
func (me *TableIO[T]) genSelectList() string {
	list := ""

	for i, j := range me.dbFields {
		list += j.FieldName
		if i < len(me.dbFields)-1 {
			list += ","
		}
	}
	return list
}

func GetTableName[T any]() string {
	var data []T

	// Tablename is dervied from the name of the objects's Type
	tableName := reflect.TypeOf(data).String()

	// Remove package names from the resulting string
	if strings.Contains(tableName, ".") {
		tableName = strings.Split(tableName, ".")[1]
	}

	// Pluralize the table name
	pluralize := pluralize.NewClient()
	tableName = pluralize.Plural(tableName)

	// Convert to snake case
	tableName = strcase.ToSnake(tableName)

	return tableName
}
