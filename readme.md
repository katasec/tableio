# Overview 

`tableio` helps with quick proto-typing in persisting structs into a database. Currently only sqlite is supported.

## Define a struct

Create a struct to represent the data you want to persist:

```go
type Hello struct {
	Message string `db:"message"`
}
```

## Connect to your DB

The TableIO constructor `NewTableIO` creates a connection to your database and returns a handle to it. Specify your struct's type as a *type parameter*, for e.g., the below example passes `[hello]` as a type parameter:

```go
helloTable, err := NewTableIO[Hello]("sqlite3", "test.db")
```

Note the database drivername and connection string are passed in the constructor.

The fields in the type parameter are used to determine the structure of your database table. For e.g., in the above case, the struct `Hello` has a field called `Messages`. As such, 
the following table will be generated:

```sql
CREATE TABLE IF NOT EXISTS Hellos (
        message VARCHAR(255) NULL
);
```




## Create the table

Call the `CreateTableIfNotExists` method to create a table for your struct:

```go
helloTable.CreateTableIfNotExists()
```

## Insert data 
To insert data, call the insert method passing in your struct

- Single Row Insert

```go
helloTable.Insert(Hello{Message: "Hi One !"})
```
- Single Row Insert
```go
	messages := []Hello{
		{Message: "Message One !"},
		{Message: "Message Two !"},
		{Message: "Message Three !"},
	}

	helloTable.InsertMany(messages)
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
