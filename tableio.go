package main

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/katasec/tableio/reflectx"
	"github.com/katasec/utils/errx"
	_ "github.com/mattn/go-sqlite3"
)

type TableIO[T any] struct {
	DB          *sqlx.DB
	tableName   string
	dbFieldsAll string
	dbFields    []reflectx.FieldInfo
}

func NewTableIO[T any](driverName string, dataSourceName string) (*TableIO[T], error) {
	db, err := sqlx.Connect(driverName, dataSourceName)
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

func (me *TableIO[T]) All() []T {
	var data []T

	sqlCmd := "select " + me.dbFieldsAll + " from " + me.tableName

	// Run select
	err := me.DB.Select(&data, sqlCmd)
	errx.PanicOnError(err)

	// return data
	return data
}

// Insert
// Inserts data into the table
// Input:
//
//	data - The data to insert
//
// Output:
//
//	nil - The data was successfully inserted
//	error - An error occurred while inserting the data
func (me *TableIO[T]) Insert(data T) error {

	sqlCmd := "insert into " + me.tableName + "(" + me.dbFieldsAll + ") values (" + reflectx.GetStructValues(data) + ")"

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

	return nil
}

func (me *TableIO[T]) DeleteTableIfExists() {

	tableName := reflectx.GetTableName[T]()

	// Start Create Table Commands
	sqlCmd := "DROP TABLE IF EXISTS " + tableName + ";"

	// Execute SQL to create table
	_, err := me.DB.Exec(sqlCmd)
	errx.PanicOnError(err)

}

func (me *TableIO[T]) Close() {
	me.DB.Close()
}
