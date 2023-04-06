package main

import "fmt"

type Shape struct {
	Width  int    `db:"width"`
	Height int    `db:"height"`
	Color  string `db:"color"`
}

type Shapex struct {
	Name       string `db:"name"`
	Dimensions dims   `db:"xx"`
}

type dims struct {
	Width  int `db:"width"`
	Height int `db:"height"`
}

func main() {
	shapexTable, _ := NewTableIO[Shapex]("sqlite3", "test.db")
	shapexTable.DeleteTableIfExists()
	shapexTable.CreateTableIfNotExists()

	shapexTable.Insert(
		Shapex{
			Name: "square",
			Dimensions: dims{
				Width:  110,
				Height: 210,
			},
		},
	)

	// Output table
	shapes := shapexTable.All()
	for i, shape := range shapes {
		fmt.Printf("%d. Color:%s, Height:%d, Width:%d \n", i, shape.Name, shape.Dimensions.Height, shape.Dimensions.Width)
	}

	shapexTable.Close()
}

func main2() {

	// Create table
	shapeTable, _ := NewTableIO[Shape]("sqlite3", "test.db")
	shapeTable.CreateTableIfNotExists()

	// Single Insert
	shapeTable.Insert(Shape{
		Width: 110, Height: 210, Color: "yellow",
	})

	// Multi Insert
	shapeTable.InsertMany([]Shape{
		{Width: 110, Height: 210, Color: "yellow"},
		{Width: 120, Height: 220, Color: "red"},
		{Width: 130, Height: 230, Color: "blue"},
	})

	// Output table
	shapes := shapeTable.All()
	for i, shape := range shapes {
		fmt.Printf("%d. Color:%s, Height:%d, Width:%d \n", i, shape.Color, shape.Height, shape.Width)
	}
}
