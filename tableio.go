package tableio

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/katasec/tableio/reflectx"
	"github.com/katasec/utils/errx"
	_ "github.com/mattn/go-sqlite3"
)

type TableIO[T any] struct {
	DB          *sql.DB
	tableName   string
	dbFieldsAll string
	dbFields    []reflectx.FieldInfo
}

func NewTableIO[T any](driverName string, dataSourceName string) (*TableIO[T], error) {
	//db, err := sqlx.Connect(driverName, dataSourceName)
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	tableio := &TableIO[T]{
		DB:        db,
		tableName: reflectx.GetTableName[T](),
		dbFields:  reflectx.GetDbStructFields[T](),
	}

	// Construct select list for table
	list := ""
	for i, j := range tableio.dbFields {
		list += j.FieldName
		if i < len(tableio.dbFields)-1 {
			list += ","
		}
	}
	tableio.dbFieldsAll = list

	return tableio, nil
}

func (me *TableIO[T]) Insert(data T) error {

	sqlCmd := "insert into " + me.tableName + "(" + me.dbFieldsAll + ") values (" + reflectx.GetStructValues(data) + ")"

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
		sqlCmd += "insert into " + me.tableName + "(" + me.dbFieldsAll + ") values (" + reflectx.GetStructValues(item) + "); \n"
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

	tableName := reflectx.GetTableName[T]()

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

	tableName := reflectx.GetTableName[T]()

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
	sqlCmd := "select " + me.dbFieldsAll + " from " + me.tableName
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

// All2 Reads all rows form the database via a "select * from <table>"
func (me *TableIO[T]) All2() []T {

	// Construct select statement
	sqlCmd := "select " + me.dbFieldsAll + " from " + me.tableName
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

		// Use reflection to get count of fields in struct
		var data T
		dataType := reflect.TypeOf(data)
		count := dataType.NumField()

		// Loop through fields and create map of field name and type
		for i := 0; i < count; i++ {
			field := dataType.Field(i)

			fieldName := field.Name
			fieldType := field.Type

			if z, ok := (scanArgs[i]).(*sql.NullBool); ok {
				masterData[fieldName] = z.Bool
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullString); ok {
				masterData[fieldName] = z.String
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt64); ok {
				masterData[fieldName] = z.Int64
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullFloat64); ok {
				masterData[fieldName] = z.Float64
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt32); ok {
				masterData[fieldName] = z.Int32
				continue
			}

			masterData[fieldName] = scanArgs[i] //reflect.New(fieldType)
			if fieldType.String() != "string" {
				masterData[fieldName], _ = strconv.Unquote(scanArgs[i].(string))
			}

		}

		finalRows = append(finalRows, masterData)
	}

	z, err := json.Marshal(finalRows)
	//fmt.Println(string(z))
	errx.PanicOnError(err)

	var data []T
	json.Unmarshal(z, &data)
	return data
}
