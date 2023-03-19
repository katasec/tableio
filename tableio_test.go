package tableio

import (
	"fmt"
	"testing"

	"github.com/katasec/utils/errx"
)

type Hello struct {
	Message string `db:"message"`
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
