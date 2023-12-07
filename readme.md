# Overview 

`tableio` is a small and simple library that helps with persisting structs in a database. Great for quick proto-typing. Only tested PostgreSQL

## Define your structs

Create structs to represent the data you want to persist. Here's an example:


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

Let's attempt to store the above `Person` struct in the database. The struct you want to persist in the DB must have a `ID` and `Name` field.

## Connect to your DB

The TableIO constructor `NewTableIO` creates a connection to your database and returns a handle to it. Specify your struct's type as a *type parameter*, for e.g., the below example passes `[Person]` as a type parameter:

### MySQL

```go
peopleTable, err := NewTableIO[Person]("mysql", "user:password@tcp(127.0.01:3306)/mydb")
```

### PostgreSQL
```go
peopleTable, err := NewTableIO[Person]("postgres", "postgresql://user:password@127.0.01/ark?sslmode=disable")
```

Note the database drivername and connection string are passed in the constructor.

The fields of your struct (that is passed via the type parameter) are used to determine the structure of your database table. For e.g., in the above case `Person` was passed, as such the fields in the `Person` struct will generate the following SQL at table creation:


```sql
CREATE TABLE IF NOT EXISTS people (
	ID SERIAL PRIMARY KEY,
	Name VARCHAR(255) NOT NULL UNIQUE,
	Age INTEGER NULL,
	Address JSONB NULL
);
```

Any complex type (i.e. not an int/string etc), is stored as a JSONB in the DB. For example, in the above case, `Address` is stored as a json object in that table. These are  transparently marshalled/unmarshalled on read/write.


## Create the table

Call the `CreateTableIfNotExists` method to create a table for your struct:

```go
peopleTable.CreateTableIfNotExists()
```

## Insert data 
To insert data, call the insert method passing in your struct

### Single Row Insert

```go
// Insert data in to table
person := Person{
	Name: "John",
	Age:  30,
	Address: Address{
		City:  "New York",
		State: "NY",
	},
}
peopleTable.Insert(person)
```
### Multiple Row Insert

```go
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
peopleTable.InsertMany(people)
```

You data is now saved in the DB !


## Query Data

### Read All data 

Call the `All()` func to retrieve all the data:

```go
data,_ := peopleTable.All()
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

### Query by Id

Call the `ById()` func to retrieve record with the given Id
```go

result, _ := peopleTable.ById(2)
person := result[0]

fmt.Printf("ID:%d Name:%s Age:%d City:%s \n", person.ID, person.Name, person.Age, person.Address.City)

```
Output:

```
ID:2 Name:Ahmed Age:45 City:Cairo 
```
### Query by Name

Call the `ByName()` func to retrieve record with the given Id
```go
result, _ := peopleTable.ByName("John")
person := result[0]

fmt.Printf("ID:%d Name:%s Age:%d City:%s \n", person.ID, person.Name, person.Age, person.Address.City)
```

Output:

```
ID:1 Name:John Age:30 City:New York 
```

## Delete Data

### Delete By Id
Call the `DeleteId()` func to delete the record with the given Id

```go
peopleTable.DeleteId(1)
```

### Delete By Name
Call the `DeleteByName()` func to delete the record with the given name:

```go
peopleTable.DeleteByName("Ahmed")
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
