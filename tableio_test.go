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
	ID      int
	Name    string
	Age     int
	Address Address
}

func TestCreateTableMySql(t *testing.T) {
	fmt.Println("Testing MySql")
	// Get connection string for env
	conn := os.Getenv("MYSQL_CONNECTION_STRING")
	if conn == "" {
		fmt.Println("Error, could not get connectin string from env var MYSQL_CONNECTION_STRING")
		os.Exit(1)
	}

	//	Create peopleTable struct from struct definition
	//	Provide DB connection string and driver name
	peopleTable, err := NewTableIO[Person]("mysql", conn)
	errx.PanicOnError(err)

	ExecTableOperations(peopleTable)
}

func TestCreateTablePgSql(t *testing.T) {
	// Get connection string for env
	conn := os.Getenv("PGSQL_CONNECTION_STRING")
	if conn == "" {
		fmt.Println("Error, could not get connectin string from env var PGSQL_CONNECTION_STRING")
		os.Exit(1)
	}

	//	Create peopleTable struct from struct definition
	//	Provide DB connection string and driver name
	peopleTable, err := NewTableIO[Person]("postgres", conn)

	//CreateTable[Person]("postgres", conn)
	errx.PanicOnError(err)

	ExecTableOperations(peopleTable)
}

// func CreateTable[T any](driverName string, datsourceName string) *TableIO[T] {
// 	// Get connection string for env
// 	conn := os.Getenv("PGSQL_CONNECTION_STRING")
// 	if conn == "" {
// 		fmt.Println("Error, could not get connectin string from env var PGSQL_CONNECTION_STRING")
// 		os.Exit(1)
// 	}

// 	//	Create peopleTable struct from struct definition
// 	//	Provide DB connection string and driver name
// 	table, err := NewTableIO[T]("postgres", conn)
// 	errx.PanicOnError(err)

//		return table
//	}

func ExecTableOperations(table *TableIO[Person]) {

	// Delete and Recreate Table
	table.DeleteTableIfExists()
	table.CreateTableIfNotExists()

	// Insert data in to table
	InsertOne(table)

	// Insert many
	InsertMany(table)

	// Read Data
	data, _ := table.All(true)
	for i, person := range data {
		fmt.Printf("%d. ID:%d Name:%s Age:%d City:%s \n", i+1, person.ID, person.Name, person.Age, person.Address.City)
	}

	// Delete table
	//table.DeleteTableIfExists()

	// Close DB Connection
	table.Close()
}

func InsertOne(table *TableIO[Person]) {
	person := Person{
		Name: "John",
		Age:  30,
		Address: Address{
			City:  "New York",
			State: "NY",
		},
	}
	table.Insert(person)
}

func InsertMany(table *TableIO[Person]) {
	people := []Person{
		{
			Name: "Ahmed",
			Age:  45,
			Address: Address{
				City:  "Cairo",
				State: "Cairo",
			},
		},
		{
			Name: "Jack",
			Age:  6,
			Address: Address{
				City:  "Terra Haute",
				State: "Indiana",
			},
		},
	}
	table.InsertMany(people)
}
