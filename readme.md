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

The TableIO constructor `NewTableIO` creates a connection to your database and returns a handle to it. Sepcify your struct's type as a type parameter, for e.g., the below exampple passed `[hello]` as a type parameter:

```go
helloTable, err := NewTableIO[Hello]("sqlite3", "test.db")
```

Note the database drivername and connection string are also passed above.


## Create the table

Call the `CreateTableIfNotExists` method to create a table for your struct:

```go
helloTable.CreateTableIfNotExists()
```

## Insert data
To insert data, call the insert method passing in your struct
```go
helloTable.Insert(Hello{Message: "Hi One !"})
helloTable.Insert(Hello{Message: "Hi Two !"})
helloTable.Insert(Hello{Message: "Hi Three !"})
```
You data is now saved in the DB !


## Read the data 

Call the `All()` method to retrieve all the data:

```go
data := helloTable.All()
for _, item := range data {
    fmt.Println(item.Message)
}
```

## Delete Table
The following deletes the table:

```go
helloTable.DeleteTableIfExists()
```

## Close DB connections

Close the DB connection

```go
// Close DB connection
helloTable.Close()
```
