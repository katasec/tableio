package main

import "fmt"

type Shape struct {
	Width  int    `db:"width"`
	Height int    `db:"height"`
	Color  string `db:"color"`
}

type Shapex struct {
	Name       string `db:"name"`
	Dimensions dims   `db:"Dimensions"`
}

type dims struct {
	Width  int `db:"Width"`
	Height int `db:"Height"`
}

func main() {
	shapexTable, _ := NewTableIO[Shapex]("sqlite3", "test.db")
	shapexTable.DeleteTableIfExists()
	shapexTable.CreateTableIfNotExists()

	shapexTable.Insert(Shapex{
		Name: "Square",
		Dimensions: dims{
			Width:  100,
			Height: 100,
		},
	})
	shapes := shapexTable.All2()

	for i, j := range shapes {
		fmt.Printf("%d. Name:%s, Height:%d, Width:%d \n", i, j.Name, j.Dimensions.Height, j.Dimensions.Width)
	}

	shapexTable.Close()
}
