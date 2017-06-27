package main

import (
	"fmt"
	"github.com/ywwg/packinglib"
)

func main() {
	t := packinglib.Trip{
		Days: 4,
		C:    &packinglib.FireflyContext,
	}

	for category, items := range t.MakeList() {
		fmt.Printf("%s:\n", category)
		for _, i := range items {
			fmt.Println("\t" + i.String())
		}
	}
}
