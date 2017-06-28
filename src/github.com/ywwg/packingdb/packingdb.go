package main

import (
	"fmt"
	"sort"

	"github.com/ywwg/packinglib"

	_ "github.com/ywwg/contexts"
	_ "github.com/ywwg/items"
)

func main() {
	t := packinglib.Trip{
		Days: 4,
		C:    packinglib.GetContext("firefly"),
	}

	packList := t.MakeList()
	// map iteration is nondeterministic so sort the keys.
	var keys []string
	for k := range packList {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, category := range keys {
		fmt.Printf("%s:\n", category)
		for _, i := range packList[category] {
			fmt.Println("\t" + i.String())
		}
	}
}
