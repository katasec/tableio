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

func TestCreateTablePgSql(t *testing.T) {

	fmt.Println("Testing PgSql")

	// Get connection string for env
	conn := os.Getenv("PGSQL_CONNECTION_STRING")
	if conn == "" {
		fmt.Println("Error, could not get connectin string from env var PGSQL_CONNECTION_STRING")
		os.Exit(1)
	}

	//	Create peopleTable struct from struct definition
	//	Provide DB connection string and driver name
	peopleTable, err := NewTableIO[Person]("postgres", conn)
	errx.PanicOnError(err)

	ExecTableOperations(peopleTable)
}

func TestPgSqlReadByName(t *testing.T) {

	fmt.Println("Reading PgSql")

	// Get connection string for env
	conn := os.Getenv("PGSQL_CONNECTION_STRING")
	if conn == "" {
		fmt.Println("Error, could not get connectin string from env var PGSQL_CONNECTION_STRING")
		os.Exit(1)
	}

	// Create peopleTable
	peopleTable, err := NewTableIO[Person]("postgres", conn)
	errx.PanicOnError(err)

	result, err := peopleTable.ByName("John")
	errx.PanicOnError(err)

	person := result[0]
	fmt.Printf("ID:%d Name:%s Age:%d City:%s \n", person.ID, person.Name, person.Age, person.Address.City)
}

func TestPgSqlReadById(t *testing.T) {

	fmt.Println("Reading PgSql")

	// Get connection string for env
	conn := os.Getenv("PGSQL_CONNECTION_STRING")
	if conn == "" {
		fmt.Println("Error, could not get connectin string from env var PGSQL_CONNECTION_STRING")
		os.Exit(1)
	}

	// Create peopleTable
	peopleTable, err := NewTableIO[Person]("postgres", conn)
	errx.PanicOnError(err)

	result, err := peopleTable.ById(2)
	errx.PanicOnError(err)

	person := result[0]
	fmt.Printf("ID:%d Name:%s Age:%d City:%s \n", person.ID, person.Name, person.Age, person.Address.City)
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

func ExecTableOperations(table *TableIO[Person]) {

	// Delete and Recreate Table
	table.DeleteTableIfExists()
	table.CreateTableIfNotExists()

	// Insert data in to table
	insertOnePerson(table)

	// Insert many
	insertMorePeople(table)

	// Read Data
	data, _ := table.All()
	for i, person := range data {
		fmt.Printf("%d. ID:%d Name:%s Age:%d City:%s \n", i+1, person.ID, person.Name, person.Age, person.Address.City)
	}

	// Delete table
	//table.DeleteTableIfExists()

	// Close DB Connection
	table.Close()
}

func TestPgSqlDeleteById(t *testing.T) {
	fmt.Println("Reading PgSql")

	// Get connection string for env
	conn := os.Getenv("PGSQL_CONNECTION_STRING")
	if conn == "" {
		fmt.Println("Error, could not get connectin string from env var PGSQL_CONNECTION_STRING")
		os.Exit(1)
	}

	// Create peopleTable
	peopleTable, err := NewTableIO[Person]("postgres", conn)
	errx.PanicOnError(err)
	peopleTable.DeleteTableIfExists()
	peopleTable.CreateTableIfNotExists()

	// Insert Data
	insertMorePeople(peopleTable)

	// Output Data
	fmt.Println("Before Delete")
	outputPeopleTable(peopleTable)

	// Delete Data
	peopleTable.DeleteId(1)

	// Output Data
	fmt.Println("After Delete")
	outputPeopleTable(peopleTable)

}

func TestPgSqlDeleteByName(t *testing.T) {

	// Get connection string for env
	conn := os.Getenv("PGSQL_CONNECTION_STRING")
	if conn == "" {
		fmt.Println("Error, could not get connectin string from env var PGSQL_CONNECTION_STRING")
		os.Exit(1)
	}

	// Create peopleTable
	peopleTable, err := NewTableIO[Person]("postgres", conn)
	peopleTable.DeleteTableIfExists()
	peopleTable.CreateTableIfNotExists()
	errx.PanicOnError(err)

	// Insert Data
	insertOnePerson(peopleTable)
	insertMorePeople(peopleTable)

	// Output Data
	fmt.Println("Before Delete")
	outputPeopleTable(peopleTable)

	// Delete Data
	peopleTable.DeleteByName("Ahmed")

	// Output Data
	fmt.Println("After Delete")
	outputPeopleTable(peopleTable)

}

func outputPeopleTable(table *TableIO[Person]) {
	// Read Data
	data, _ := table.All()
	for i, person := range data {
		fmt.Printf("%d. ID:%d Name:%s Age:%d City:%s \n", i+1, person.ID, person.Name, person.Age, person.Address.City)
	}
}
func insertOnePerson(table *TableIO[Person]) {
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

func insertMorePeople(table *TableIO[Person]) {
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
