package main

// import "fmt"

// func main2() {

// 	// Create table
// 	shapeTable, _ := NewTableIO[Shape]("sqlite3", "test.db")
// 	shapeTable.CreateTableIfNotExists()

// 	// Single Insert
// 	shapeTable.Insert(Shape{
// 		Width: 110, Height: 210, Color: "yellow",
// 	})

// 	// Multi Insert
// 	shapeTable.InsertMany([]Shape{
// 		{Width: 110, Height: 210, Color: "yellow"},
// 		{Width: 120, Height: 220, Color: "red"},
// 		{Width: 130, Height: 230, Color: "blue"},
// 	})

// 	// Output table
// 	shapes := shapeTable.All()
// 	for i, shape := range shapes {
// 		fmt.Printf("%d. Color:%s, Height:%d, Width:%d \n", i, shape.Color, shape.Height, shape.Width)
// 	}
// }
