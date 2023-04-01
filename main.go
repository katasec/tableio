package main

import "fmt"

type Shape struct {
	Width  int    `db:"width"`
	Height int    `db:"height"`
	Color  string `db:"color"`
}

func main() {

	shapeTable, _ := NewTableIO[Shape]("sqlite3", "test.db")
	shapeTable.CreateTableIfNotExists()

	shapeTable.Insert(Shape{Width: 10, Height: 20, Color: "red"})

	shapes := shapeTable.All()

	for i, shape := range shapes {
		fmt.Printf("%d. Color:%s, Height:%d, Width:%d \n", i, shape.Color, shape.Height, shape.Width)
	}

}
