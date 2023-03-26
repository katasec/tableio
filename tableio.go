package tableio

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/katasec/tableio/reflectx"
	"github.com/katasec/utils/errx"
	_ "github.com/mattn/go-sqlite3"
)

type TableIO[T any] struct {
	DB              *sqlx.DB
	tableName       string
	dbSelectListAll string
	fields          []reflectx.FieldInfo
	dbFields        []string
}

func NewTableIO[T any](driverName string, dataSourceName string) (*TableIO[T], error) {
	db, err := sqlx.Connect(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	//selectFields := reflectx.GetDbStructFields[T]()
	//selectList := strings.Join(selectFields, ", ")

	tableio := &TableIO[T]{
		DB: db,
		//dbSelectListAll: selectList, //GetDbColumnNames[T](),
		tableName: reflectx.GetTableName[T](),
		fields:    reflectx.GetStructFieldsX[T](),
		dbFields:  reflectx.GetDbStructFields[T](),
	}

	return tableio, nil
}

func (me *TableIO[T]) SelectList() string {
	var sb strings.Builder

	for i, field := range me.fields {
		sb.WriteString(field.FieldName)
		if i < len(me.fields) {
			sb.WriteString(",")
		}
	}

	return sb.String()
}

func (me *TableIO[T]) All() []T {
	var data []T

	sqlCmd := "select " + me.dbSelectListAll + " from " + me.tableName

	// Run select
	err := me.DB.Select(&data, sqlCmd)
	errx.PanicOnError(err)

	// return data
	return data
}

func (me *TableIO[T]) Insert(data T) error {

	sqlCmd := "insert into " + me.tableName + "(" + me.dbSelectListAll + ") values (" + GetStructValues(data) + ")"

	// Run Insert
	_, err := me.DB.Exec(sqlCmd)
	errx.PanicOnError(err)

	// return data
	return nil
}

func (me *TableIO[T]) CreateTableIfNotExists() error {

	var sb strings.Builder

	tableName := reflectx.GetTableName[T]()

	// Start Create Table Commands
	sb.WriteString("CREATE TABLE IF NOT EXISTS " + tableName + " (\n")

	// Add fields
	sb.WriteString(GenSqlForStructFields[T]())

	// End Command
	sb.WriteString(");")

	// Generate string
	sqlCmd := sb.String()
	fmt.Println(sqlCmd)

	// Execute SQL to create table
	result, err := me.DB.Exec(sqlCmd)
	fmt.Println(result)
	errx.PanicOnError(err)

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
