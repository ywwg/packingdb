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

	for _, p := range t.MakeList() {
		fmt.Println(p.String())
	}
}
