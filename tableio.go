package tableio

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	// DB Drivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	// Pluralize & Snake Case
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"

	// Internal libs
	"github.com/katasec/tableio/reflectx"

	// Other libs
	"github.com/katasec/utils/errx"
)

type TableIO[T any] struct {
	DB         *sql.DB
	tableName  string
	driverName string
	// This is a comma separated list of fields for the table used for SELECT statements
	selectList string

	// This is a comma separated list of fields for the table used for INSERT statements
	// this excludes the ID field
	insertList string

	// dbFields is a list of fields in the struct that have a 'db' tag
	dbFields []reflectx.FieldInfo
}

// NewTableIO Creates a new TableIO object for a given struct
func NewTableIO[T any](driverName string, dataSourceName string) (*TableIO[T], error) {

	// Validate Struct first. A struct must have an 'ID' and 'Name' field
	if !Validate[T]() {
		return nil, fmt.Errorf("error: TableIO structs must have an 'ID' and 'Name' field")
	}

	// Create a DB connection
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	// Create TableIO struct with a connection to the DB
	tableio := &TableIO[T]{
		DB:         db,
		tableName:  GenTableName[T](),
		dbFields:   reflectx.GetStructFields[T](), //reflectx.GetDbStructFields[T](),
		driverName: driverName,
	}

	// Generate and cache 'select list' and 'insert list' for table
	tableio.selectList = tableio.genSelectList()
	tableio.insertList = tableio.genInsertList()

	return tableio, nil
}

// Validate Checks if a struct has an 'ID' and 'Name' field
func Validate[T any]() bool {

	var x T
	typeName := reflect.TypeOf(x).String()

	// Get fields in struct
	fields := reflectx.GetStructFields[T]()

	// Check if struct has an 'ID' and 'Name' field
	idField := false
	nameField := false
	for _, field := range fields {
		if field.FieldName == "ID" && field.FieldType == "int64" {
			idField = true
		}
		if field.FieldName == "Name" && field.FieldType == "string" {
			nameField = true
		}
	}

	// Output error if struct does not have an 'ID' and 'Name' field
	if !idField {
		fmt.Printf("struct '%s' does not have an 'ID' field of type int64\n", typeName)
	}
	if !nameField {
		fmt.Printf("struct '%s' does not have a 'Name' field of type string\n", typeName)
	}

	// Return true if struct has both ID and Name fields
	return idField && nameField
}

// GenTableName Generates the DB table name from the type of the struct.
// For e.g. if the struct is named 'User', the table name in the DB will be 'users'
func GenTableName[T any]() string {
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

// Insert Inserts a single row into the table
func (me *TableIO[T]) Insert(data T, verbose ...bool) error {
	var debug bool

	if len(verbose) > 0 {
		debug = verbose[0]
	}

	sqlCmd := "insert into " + me.tableName + "(" + me.insertList + ") values (" + reflectx.GetStructValuesForInsert(data) + ")"

	if debug {
		fmt.Println(sqlCmd)
	}

	// Run Insert
	_, err := me.DB.Exec(sqlCmd)
	errx.PanicOnError(err)

	// return data
	return nil
}

// InsertMany Inserts multiple rows into the table
func (me *TableIO[T]) InsertMany(data []T) error {

	// Gen Sql Command
	sqlCmd := ""
	for _, item := range data {
		sqlCmd += "insert into " + me.tableName + "(" + me.insertList + ") values (" + reflectx.GetStructValuesForInsert(item) + "); \n"
	}

	// Run Insert
	_, err := me.DB.Exec(sqlCmd)
	errx.PanicOnError(err)

	// return data
	return nil
}

// CreateTable Creates a table in the DB for the struct if it does not exist
func (me *TableIO[T]) CreateTableIfNotExists(verbose ...bool) error {
	var debug bool

	if len(verbose) > 0 {
		debug = verbose[0]
	}

	var sb strings.Builder

	// Generate a table name to use for creating in the DB
	tableName := GenTableName[T]()

	// Start Create Table Command
	sb.WriteString("CREATE TABLE IF NOT EXISTS " + tableName + " (\n")

	// Add fields
	sb.WriteString(reflectx.GenSqlForFields(me.dbFields, me.driverName))

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

	if debug {
		fmt.Println("Create table '" + me.tableName + "' successfully")
	}

	return nil
}

// DeleteTableIfExists Deletes a table in the DB for the struct if it exists
func (me *TableIO[T]) DeleteTableIfExists(verbose ...bool) {

	// Get verbose flag
	var debug bool
	if len(verbose) > 0 {
		debug = verbose[0]
	}

	tableName := GenTableName[T]()

	// Start Create Table Commands
	sqlCmd := "DROP TABLE IF EXISTS " + tableName + ";"

	// Execute SQL to create table
	_, err := me.DB.Exec(sqlCmd)
	errx.PanicOnError(err)

	// Print SQL if debug flag is set
	if debug {
		fmt.Println("Deleted table '" + me.tableName + "' successfully.")
		fmt.Println(sqlCmd)
	}
}

// Close Closes the DB connection
func (me *TableIO[T]) Close() {
	me.DB.Close()
}

// All Returns all rows in the table
func (me *TableIO[T]) All(verbose ...bool) []T {

	var debug bool

	if len(verbose) > 0 {
		debug = verbose[0]
	}

	// Construct select statement
	sqlCmd := "select " + me.selectList + " from " + me.tableName
	if debug {
		fmt.Println(sqlCmd)
	}

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

		// Create map of column names and values of appropriate type
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

		// Append row to final result
		finalRows = append(finalRows, masterData)
	}

	// Marshal final result to JSON
	jsonString, err := json.Marshal(finalRows)
	fmt.Println(string(jsonString))
	errx.PanicOnError(err)

	// Unmarshal JSON to struct
	var data []T
	json.Unmarshal(jsonString, &data)

	// Return data of type T
	return data
}

// genSelectList returns a comma separated list of fields for the table used for select statements
func (me *TableIO[T]) genSelectList() string {
	list := ""

	// Loop through fields
	for i, j := range me.dbFields {

		// Add field name to list
		list += j.FieldName

		// Add comma if not last field
		if i < len(me.dbFields)-1 {
			list += ","
		}
	}
	return list
}

// genSelectList returns a comma separated list of fields for the table used for insert statements
// this excludes the ID field
func (me *TableIO[T]) genInsertList() string {
	list := ""

	for i, j := range me.dbFields {

		// Add field name to list
		if j.FieldName != "ID" {
			list += j.FieldName

			// Add comma if not last field
			if i < len(me.dbFields)-1 {
				list += ","
			}
		}

	}
	return list
}
