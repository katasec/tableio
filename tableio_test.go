package tableio

import (
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
	helloTable.Insert(Hello{
		Message: "Hi there !",
	})

	// Delete table
	helloTable.DeleteTableIfExists()

	// Close DB connection
	helloTable.Close()
}
