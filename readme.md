# Overview 

`tableio` helps with quick proto-typing in persisting structs into a database. Only tested with Mysql and PostgreSQL

## Define your structs

Create structs to represent the data you want to persist:

```go
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
```

Note in this example we'll attempt to store a `Person` in the database. The struct you want to persist in the DB must have a `ID` and `Name` field.

## Connect to your DB

The TableIO constructor `NewTableIO` creates a connection to your database and returns a handle to it. Specify your struct's type as a *type parameter*, for e.g., the below example passes `[hello]` as a type parameter:

### MySQL:

```go
peopleTable, err := NewTableIO[Hello]("mysql", "user:password@tcp(127.0.01:3306)/mydb")
```

### PostgreSQL
```go
peopleTable, err := NewTableIO[Shapex]("postgres", "postgresql://user:password@127.0.01/ark?sslmode=disable")
```

Note the database drivername and connection string are passed in the constructor.

The fields in the type parameter are used to determine the structure of your database table. For e.g., in the above case, the fields in the `Address` and `Person` struct will generate the following SQL for table creation:


```sql
CREATE TABLE IF NOT EXISTS people (
	ID SERIAL PRIMARY KEY,
	Name VARCHAR(255) NOT NULL UNIQUE,
	Age INTEGER NULL,
	Address JSONB NULL
);
```

Any complex type (i.e. not an int/string etc), is stored as a JSONB in the DB. It's automatically marshalled/unmarshalled on read/write


## Create the table

Call the `CreateTableIfNotExists` method to create a table for your struct:

```go
peopleTable.CreateTableIfNotExists()
```

## Insert data 
To insert data, call the insert method passing in your struct

- Single Row Insert

```go
	peopleTable.Insert(Person{
		Name: "John",
		Age:  30,
		Address: Address{
			City:  "New York",
			State: "NY",
		},
	})
```
- Multiple Row Insert
```go
	peopleTable.InsertMany(
		[]Person{
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
		},
	)
```

You data is now saved in the DB !


## Read the data 

Call the `All()` method to retrieve all the data:

```go
	data := peopleTable.All()
	for i, person := range data {
		fmt.Printf("%d. ID:%d Name:%s Age:%d City:%s \n", i+1, person.ID, person.Name, person.Age, person.Address.City)
	}
```

Output:
```
1. ID:1 Name:John Age:30 City:New York
2. ID:2 Name:Ahmed Age:45 City:Cairo
3. ID:3 Name:Jack Age:6 City:Terra Haute

```
## Delete Table
The following deletes the table:

```go
peopleTable.DeleteTableIfExists()
```

## Close DB connection

Close the DB connection

```go
// Close DB connection
peopleTable.Close()
```
