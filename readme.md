# Overview 

TableIO helps with persisting structs into your datbase. Currently only sqlite is supported.

## Define a struct

Create a struct to represent the data you want to persist:

```go
type Hello struct {
	Message string `db:"message"`
}
```

## Connect to your DB

The TableIO constructor `NewTableIO` creates a connection to your database and returns a handle to it. Pass in your database driver name and connection string:

```go
helloTable, err := NewTableIO[Hello]("sqlite3", "test.db")
```

## Create the table

Call the `CreateTableIfNotExists` method to create a table for your struct:

```go
helloTable.CreateTableIfNotExists()
```

## Insert data
To insert data, call the insert method passing in your struct
```go
helloTable.Insert(Hello{
    Message: "Hi there !",
})
```
You data is now saved in the DB !

## Delete Table
The following deletes the table:

```go
helloTable.DeleteTableIfExists()
```

## Close DB connections

Close the DB connection

```
// Close DB connection
helloTable.Close()
```
