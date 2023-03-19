# Overview 

TableIO helps with persisting structs to sqlite. For e.g.

## Define a struct

```go
type Hello struct {
	Message string `db:"message"`
}
```

## Create a table using the struct's definition

The TableIO constructor `NewTableIO` creates the table in your database and returns a handle to it. Pass in your struct's type as a type parameter into the method as per below:

```go
helloTable, err := NewTableIO[Hello]("sqlite3", "test.db")
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
