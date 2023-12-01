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
	if err != nil {
		message := fmt.Sprintf("Error creating table %s: %s", tableName, err.Error())
		return errors.New(message)
	} else {
		if debug {
			fmt.Println("Create table '" + me.tableName + "' successfully")
		}
		return nil
	}
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

// Close Closes the DB connection
func (me *TableIO[T]) Close() {
	me.DB.Close()
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

// All Returns all rows in the table
func (me *TableIO[T]) All(verbose ...bool) ([]T, error) {

	// Configure verbose flag
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
	if err != nil {
		return nil, err
	}

	// Get column types and count
	dbColumnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	// Loop through rows
	allRows := []interface{}{} // This will contain the rows returned from the DB
	for rows.Next() {
		// Create array of pointers for each column of appropriate type
		dest := me.initializeScanDest(dbColumnTypes)

		// Scan row into dest
		err := rows.Scan(dest...)
		if err != nil {
			return nil, err
		}

		// Dest is a []interface{}. Casr the items in the array to the appropriate type
		// and store them in the map row
		row := me.castColumnTypes(dest, dbColumnTypes)

		// Append current row to final result
		allRows = append(allRows, row)
	}

	// Marshal final result to JSON
	jsonBytes, err := json.MarshalIndent(allRows, "", "  ")
	if err != nil {
		return nil, err
	}

	// For some reason the jsonBytes needs to be converted to a string and then back to bytes
	// to cast it to a []T

	// Convert jsonBytes to string
	jsonString := string(jsonBytes)

	// The resulting JSON string has escaped quotes and curly braces. Unescape them
	jsonString = UnescapeJson(jsonString)

	// Convert jsonString to []T
	var data []T
	json.Unmarshal([]byte(jsonString), &data)

	// Return data
	return data, nil

}

// UnescapeJson Unescapes a JSON string
func UnescapeJson(jsonString string) string {

	// Fix escaped quotes
	jsonString = strings.Replace(jsonString, `\"`, `"`, -1)

	// Fix escaped curly braces
	jsonString = strings.Replace(jsonString, `"{`, `{`, -1)
	jsonString = strings.Replace(jsonString, `}"`, `}`, -1)

	return jsonString
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
