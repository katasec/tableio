package tableio

import (
	"fmt"
	"os"
	"testing"

	// DB Drivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/katasec/utils/errx"
)

// Create a test struct

type Address struct {
	City  string
	State string
}
type Person struct {
	ID      int64
	Name    string
	Age     int
	Address Address
}

func TestCreateTableMySql(t *testing.T) {
	// Get connection string for env
	conn := os.Getenv("MYSQL_CONNECTION_STRING")
	if conn == "" {
		fmt.Println("Error, could not get connectin string from env var MYSQL_CONNECTION_STRING")
		os.Exit(1)
	}

	/*
		Create peopleTable struct from struct definition
		Provide DB connection string and driver name
	*/
	peopleTable, err := NewTableIO[Person]("mysql", conn)
	errx.PanicOnError(err)

	// Delete and Recreate Table
	peopleTable.DeleteTableIfExists(true)
	peopleTable.CreateTableIfNotExists(true)

	// Insert data in to table
	peopleTable.Insert(Person{
		Name: "John",
		Age:  30,
		Address: Address{
			City:  "New York",
			State: "NY",
		},
	})
	peopleTable.Insert(Person{
		Name: "Ahmed",
		Age:  45,
		Address: Address{
			City:  "Abu Dhabi",
			State: "UAE",
		},
	})

	// Read Data
	data := peopleTable.All()
	for _, person := range data {
		fmt.Println(person.Name)
	}
	// Delete table
	//peopleTable.DeleteTableIfExists()

	// Close DB connection
	peopleTable.Close()
}

//	func TestTable(table *TableIO[Person]) {
//		// Insert data in to table
//		table.Insert(Person{
//			Name: "John",
//			Age:  30,
//			Address: Address{
//				City:  "New York",
//				State: "NY",
//			},
//		})
//		table.Insert(Person{
//			Name: "Ahmed",
//			Age:  45,
//			Address: Address{
//				City:  "Abu Dhabi",
//				State: "UAE",
//			},
//		})
//	}
func TestCreateTablePgSql(t *testing.T) {
	// Get connection string for env
	conn := os.Getenv("PGSQL_CONNECTION_STRING")
	if conn == "" {
		fmt.Println("Error, could not get connectin string from env var PGSQL_CONNECTION_STRING")
		os.Exit(1)
	}

	/*
		Create peopleTable struct from struct definition
		Provide DB connection string and driver name
	*/
	peopleTable, err := NewTableIO[Person]("postgres", conn)
	errx.PanicOnError(err)

	// Delete and Recreate Table
	peopleTable.DeleteTableIfExists(true)
	peopleTable.CreateTableIfNotExists(true)

	// Insert data in to table
	peopleTable.Insert(Person{
		Name: "John",
		Age:  30,
		Address: Address{
			City:  "New York",
			State: "NY",
		},
	})
	peopleTable.Insert(Person{
		Name: "Ahmed",
		Age:  45,
		Address: Address{
			City:  "Abu Dhabi",
			State: "UAE",
		},
	})

	// Read Data
	data := peopleTable.All()
	for _, person := range data {
		fmt.Println(person.Name)
	}
	// Delete table
	//peopleTable.DeleteTableIfExists()

	// Close DB connection
	peopleTable.Close()
}

// func TestGenSqlForFields(t *testing.T) {

// 	fields := reflectx.GetDbStructFields[Hello]()

// 	x := reflectx.GenSqlForFields(fields)

// 	fmt.Println(x)
// }

// func TestSelectList(t *testing.T) {
// 	//helloTable, _ := NewTableIO[Entity]("sqlite3", "test.db")

// 	x := reflectx.GetStructFields[Entity]()

// 	for i, field := range x {
// 		fmt.Println(i, "Name:"+field.FieldName, "Type:"+field.FieldType)
// 	}
// 	//fmt.Println("Fields: " + helloTable.selectList)
// }

// type AzureCloudspace struct{}

// func TestTableNaming(t *testing.T) {
// 	fmt.Println(GenTableName[Entity]())
// 	fmt.Println(GenTableName[AzureCloudspace]())
// }

func TestValidateStruct(t *testing.T) {
	type TestStruct1 struct {
		Age int
	}

	_, err := NewTableIO[TestStruct1]("", "")
	if err != nil {
		fmt.Println(err)
	}

	type TestStruct2 struct {
		ID   int64
		Name string
		Age  int
	}

	_, err = NewTableIO[TestStruct2]("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("TestStruct2 is a valid TableIO struct")
	}
}
