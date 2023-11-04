package main

import (
	"fmt"
)

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

	// Create new table of type [shapex] for mysql
	shapexTable, _ := NewTableIO[Shapex]("mysql", "root:Password123@tcp(127.0.01:3306)/ark")
	shapexTable.DeleteTableIfExists()
	//os.Exit(0)
	shapexTable.CreateTableIfNotExists()

	shapexTable.Insert(Shapex{
		Name: "square1",
		Dimensions: dims{
			Width:  100,
			Height: 100,
		},
	})

	shapexTable.Insert(Shapex{
		Name: "square2",
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
