package main

import (
	"fmt"
	"github.com/ywwg/packinglib"

	_ "github.com/ywwg/items"
)

func main() {
	t := packinglib.Trip{
		Days: 4,
		C:    packinglib.GetContext("firefly"),
	}

	for category, items := range t.MakeList() {
		fmt.Printf("%s:\n", category)
		for _, i := range items {
			fmt.Println("\t" + i.String())
		}
	}
}
