package tableio

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/katasec/tableio/reflectx"
	"github.com/katasec/utils/errx"
)

type Entity struct {
	Id      int
	Name    string
	Entity2 Entity2
	Field1  string
	Field2  string
}

type Entity2 struct {
	Id     int
	Name   string
	Field1 string
	Field2 string
}

// Create a test struct
type Hello struct {
	Message  string `db:"message"`
	Message1 string `db:"message1"`
	Message2 string `db:"message2"`
}

func TestStuff(t *testing.T) {
	x := Entity{
		Id:   1,
		Name: "name",
		Entity2: Entity2{
			Id:     2,
			Name:   "name2",
			Field1: "field1",
			Field2: "field2",
		},
		Field1: "field1",
	}

	xBytes, _ := json.Marshal(x)
	fmt.Println(string(xBytes))
}
func TestCreateTable(t *testing.T) {

	// Get connection string for env
	conn := os.Getenv("MYSQL_CONNECTION_STRING")
	if conn == "" {
		fmt.Println("Error, could not get connectin string from env var MYSQL_CONNECTION_STRING")
		os.Exit(1)
	}
	// Create New Table from struct definition
	helloTable, err := NewTableIO[Hello]("mysql", conn)
	errx.PanicOnError(err)

	helloTable.CreateTableIfNotExists(true)

	// Insert data in to table
	helloTable.Insert(Hello{
		Message:  "message",
		Message1: "message1",
		Message2: "message2",
	})
	helloTable.Insert(Hello{
		Message:  "2message",
		Message1: "2message1",
		Message2: "2message2",
	})

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

func TestGenSqlForFields(t *testing.T) {

	fields := reflectx.GetDbStructFields[Hello]()

	x := reflectx.GenSqlForFields(fields)

	fmt.Println(x)
}

func TestSelectList(t *testing.T) {
	//helloTable, _ := NewTableIO[Entity]("sqlite3", "test.db")

	x := reflectx.GetStructFields[Entity]()

	for i, field := range x {
		fmt.Println(i, "Name:"+field.FieldName, "Type:"+field.FieldType)
	}
	//fmt.Println("Fields: " + helloTable.selectList)
}

type AzureCloudspace struct{}

func TestTableNaming(t *testing.T) {
	fmt.Println(GetTableName[Entity]())
	fmt.Println(GetTableName[AzureCloudspace]())
}
