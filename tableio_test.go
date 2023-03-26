package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/katasec/tableio/reflectx"
	"github.com/katasec/utils/errx"
)

// Create a test struct
type Hello struct {
	Message  string `db:"message"`
	Message1 string `db:"message1"`
	Message2 string `db:"message2"`
}

func TestCreateTable(t *testing.T) {

	// Create New Table from struct definition
	helloTable, err := NewTableIO[Hello]("sqlite3", "test.db")
	errx.PanicOnError(err)

	helloTable.CreateTableIfNotExists()

	// Insert data in to table
	helloTable.Insert(Hello{Message: "Hi One !"})
	helloTable.Insert(Hello{Message: "Hi Two !"})
	helloTable.Insert(Hello{Message: "Hi Three !"})

	// Read Data
	data := helloTable.All()
	for _, item := range data {
		fmt.Println(item.Message)
	}
	// Delete table
	helloTable.DeleteTableIfExists()

	// Close DB connection
	helloTable.Close()
}

func TestGetDbStructFieldsByTag(t *testing.T) {

	// Create a test struct
	type Hello struct {
		Message0 string `db:"message0"`
		Message1 string `db:"message1"`
		Message2 string `db:"message2"`
		Message3 string
	}

	y := reflectx.GetDbStructFields[Hello]()

	fmt.Print("\nDb Select List: " + strings.Join(y, ", ") + "\n\n")

}

func TestGenSqlForFields(t *testing.T) {

	fields := reflectx.GetStructFieldsX[Hello]()

	x := reflectx.GenSqlForFields(fields)

	fmt.Println(x)
}

func TestSelectList(t *testing.T) {
	helloTable, _ := NewTableIO[Hello]("sqlite3", "test.db")

	fmt.Println(helloTable.SelectList())
	fmt.Println(helloTable.dbFieldsAll)
}
