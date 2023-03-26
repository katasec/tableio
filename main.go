package main

import "fmt"

type Shape struct {
	Width  int    `db:"width"`
	Height int    `db:"height"`
	Color  string `db:"color"`
}

func main() {

	shapeTable, _ := NewTableIO[Shape]("sqlite3", "test.db")
	//shapeTable.CreateTableIfNotExists()

	fmt.Println("Select List:" + shapeTable.dbFieldsAll)

	shapeTable.Insert(Shape{Width: 10, Height: 20, Color: "red"})

}
