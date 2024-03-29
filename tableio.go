package tableio

import (
	"database/sql"
	"encoding/json"
	"errors"
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

	// dbColumnTypes is a list of SQL column types in the DB. This is used
	// to cast the values returned from the DB to the appropriate Go field type
	dbColumnTypes []*sql.ColumnType
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

	err = db.Ping()
	if err != nil {
		message := fmt.Sprintf("Error connecting to DB: %s", err.Error())
		return nil, errors.New(message)
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

	// If table exists, cache column types
	tableio.dbColumnTypes, _ = tableio.getDbColumnTypes()

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
		if field.FieldName == "ID" && (field.FieldType == "int64" || field.FieldType == "int32" || field.FieldType == "int") {
			idField = true
		}
		if field.FieldName == "Name" && field.FieldType == "string" {
			nameField = true
		}
	}

	// Output error if struct does not have an 'ID' and 'Name' field
	if !idField {
		fmt.Printf("struct '%s' does not have an 'ID' field of type int64/int32/int\n", typeName)
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

	//Execute SQL to create table, return error if any
	_, err := me.DB.Exec(sqlCmd)
	if err != nil {
		message := fmt.Sprintf("Error creating table %s: %s", tableName, err.Error())
		return errors.New(message)
	}

	// If table was created, cache column types
	me.dbColumnTypes, err = me.getDbColumnTypes()
	if err != nil {
		return err
	}

	if debug {
		fmt.Println("Create table '" + me.tableName + "' successfully")
	}
	return nil

}

// DeleteTableIfExists Deletes a table in the DB for the struct if it exists
func (me *TableIO[T]) DeleteTableIfExists(verbose ...bool) error {

	// Get verbose flag
	var debug bool
	if len(verbose) > 0 {
		debug = verbose[0]
	}

	tableName := GenTableName[T]()

	// Start Create Table Commands
	sqlCmd := "DROP TABLE IF EXISTS " + tableName + ";"
	if debug {
		fmt.Println(sqlCmd)
	}
	// Execute SQL to create table
	_, err := me.DB.Exec(sqlCmd)
	if err != nil {
		message := fmt.Sprintf("Error creating table %s: %s", tableName, err.Error())
		return errors.New(message)
	} else {
		if debug {
			fmt.Println("Deleted table '" + me.tableName + "' successfully.")
		}
		return nil
	}
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
	return err
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
	return err
}

// All Returns all rows in the table
func (me *TableIO[T]) All(verbose ...bool) ([]T, error) {
	// Configure verbose flag
	debug := len(verbose) > 0 && verbose[0]

	// Construct select statement
	sqlCmd := "select " + me.selectList + " from " + me.tableName
	if debug {
		fmt.Println(sqlCmd)
	}

	return me.executeQuery(sqlCmd, me.dbColumnTypes)
}

// ByName Returns all rows in the table with the given name
func (me *TableIO[T]) ByName(name string, verbose ...bool) ([]T, error) {
	// Configure verbose flag
	debug := len(verbose) > 0 && verbose[0]

	// Construct select statement
	sqlCmd := fmt.Sprintf("select %s from %s where name = '%s'", me.selectList, me.tableName, name)
	if debug {
		fmt.Println(sqlCmd)
	}

	return me.executeQuery(sqlCmd, me.dbColumnTypes)
}

// ById Returns all rows in the table with the given ID
func (me *TableIO[T]) ById(id int, verbose ...bool) ([]T, error) {
	// Configure verbose flag
	debug := len(verbose) > 0 && verbose[0]

	// Construct select statement
	sqlCmd := fmt.Sprintf("select %s from %s where id = '%d'", me.selectList, me.tableName, id)
	if debug {
		fmt.Println(sqlCmd)
	}

	return me.executeQuery(sqlCmd, me.dbColumnTypes)
}

// Close Closes the DB connection
func (me *TableIO[T]) Close() {
	me.DB.Close()
}

// DeleteId Deletes a row with the given ID
func (me *TableIO[T]) DeleteId(id int) error {
	sqlCmd := fmt.Sprintf("delete from %s where id = '%d'", me.tableName, id)
	_, err := me.DB.Exec(sqlCmd)
	return err
}

// DeleteByName Deletes a row with the given name
func (me *TableIO[T]) DeleteByName(name string) error {
	sqlCmd := fmt.Sprintf("delete from %s where name = '%s'", me.tableName, name)
	_, err := me.DB.Exec(sqlCmd)
	return err
}

// initializeScanDest Creates an array items for each column in a DB row with the appropriate column type.
// This will be used to store the incoming data from the DB row
func (me *TableIO[T]) initializeScanDest(dbColumnTypes []*sql.ColumnType) []any {
	count := len(dbColumnTypes)

	dest := make([]any, count)
	for i, v := range dbColumnTypes {
		switch v.DatabaseTypeName() {
		case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
			dest[i] = new(sql.NullString)
		case "BOOL":
			dest[i] = new(sql.NullBool)
		case "INT4":
			dest[i] = new(sql.NullInt64)
		case "JSON":
			dest[i] = new(sql.NullString)
		case "JSONB":
			dest[i] = new(sql.NullString)
		default:
			dest[i] = new(sql.NullString)
		}
	}

	return dest
}

// castColumnTypes Casts the values in the dest array to the appropriate type and stores them in the currentRowData map
func (me *TableIO[T]) castColumnTypes(dest []any, dbColumnTypes []*sql.ColumnType) map[string]any {
	//currentRowData := map[string]interface{}{}
	currentRowData := map[string]any{}
	for i, v := range dbColumnTypes {

		if z, ok := (dest[i]).(*sql.NullBool); ok {
			currentRowData[v.Name()] = z.Bool
			continue
		}

		if z, ok := (dest[i]).(*sql.NullString); ok {
			currentRowData[v.Name()] = z.String
			continue
		}

		if z, ok := (dest[i]).(*sql.NullInt64); ok {
			currentRowData[v.Name()] = z.Int64
			continue
		}

		if z, ok := (dest[i]).(*sql.NullFloat64); ok {
			currentRowData[v.Name()] = z.Float64
			continue
		}

		if z, ok := (dest[i]).(*sql.NullInt32); ok {
			currentRowData[v.Name()] = z.Int32
			continue
		}

		if z, ok := (dest[i]).(*sql.NullByte); ok {
			currentRowData[v.Name()] = z.Byte
			continue
		}

		currentRowData[v.Name()] = dest[i]
	}

	return currentRowData
}

// Common function to handle query execution and result parsing
func (me *TableIO[T]) executeQuery(sqlCmd string, dbColumnTypes []*sql.ColumnType) ([]T, error) {
	// Run Query
	rows, err := me.DB.Query(sqlCmd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create an empty array of pointers for each column of appropriate type
	dest := me.initializeScanDest(dbColumnTypes)

	// Loop through rows
	allRows := []interface{}{}
	for rows.Next() {
		// Scan row into dest
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}

		// Cast items in the array to the appropriate type and store them in the map row
		row := me.castColumnTypes(dest, dbColumnTypes)

		// Append current row to final result
		allRows = append(allRows, row)
	}

	// Marshal final result to JSON
	jsonBytes, err := json.MarshalIndent(allRows, "", "  ")
	if err != nil {
		return nil, err
	}

	// Convert jsonBytes to string. For some reason the jsonBytes
	// needs to be converted to a string and then back to bytes
	// vs. directly casting it to a []T
	jsonString := string(jsonBytes)

	// The resulting JSON string has escaped quotes and curly braces. Unescape them
	jsonString = unescapeJson(jsonString)

	// Convert jsonString to bytes and unmarshall to []T
	var data []T
	if err := json.Unmarshal([]byte(jsonString), &data); err != nil {
		return nil, err
	}

	// Return data
	return data, nil
}

// unescapeJson Unescapes a JSON string
func unescapeJson(jsonString string) string {

	// Fix escaped quotes
	jsonString = strings.Replace(jsonString, `\"`, `"`, -1)

	// Fix escaped curly braces
	jsonString = strings.Replace(jsonString, `"{`, `{`, -1)
	jsonString = strings.Replace(jsonString, `}"`, `}`, -1)

	return jsonString
}

// getDbColumnTypes Returns a list of column types in the DB
func (me *TableIO[T]) getDbColumnTypes() ([]*sql.ColumnType, error) {

	// Execute a query that will return no rows
	sqlCmd := "SELECT * FROM " + me.tableName + " WHERE 1=0"
	rows, err := me.DB.Query(sqlCmd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get column types from the result
	dbColumnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	// Return column types
	return dbColumnTypes, nil
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
