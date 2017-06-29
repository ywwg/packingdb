package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ywwg/packinglib"

	_ "github.com/ywwg/contexts"
	_ "github.com/ywwg/items"
)

var (
	flagContext     = flag.String("context", "", "The context you want to load (not needed if packing file exists)")
	flagDays        = flag.Int("days", 0, "The number of days for the trip (not needed if packing file exists)")
	flagPackingFile = flag.String("packfile", "", "The filename to create or load (not needed if you just want to print a list)")
	flagPackItem    = flag.String("pack", "", "The name of an item that has been packed (noop without packfile)")
)

func main() {
	flag.Parse()

	var t *packinglib.Trip

	// Simple mode: just print the list.
	if *flagPackingFile == "" {
		if len(*flagContext) == 0 {
			panic("Need a context")
		}
		if *flagDays == 0 {
			panic("Need a number of days")
		}
		t = packinglib.NewTrip(*flagDays, *flagContext)
		for _, l := range t.Strings() {
			fmt.Println(l)
		}
		return
	}

	// File mode: load if the file already exists or create new if not
	if _, err := os.Stat(*flagPackingFile); !os.IsNotExist(err) {
		t = &packinglib.Trip{}
		if err2 := t.LoadFromFile(*flagPackingFile); err2 != nil {
			panic(fmt.Sprintf("%v", err2))
		}
	} else {
		t = packinglib.NewTrip(*flagDays, *flagContext)
	}

	if *flagPackItem != "" {
		t.Pack(*flagPackItem)
	}
	for _, l := range t.Strings() {
		fmt.Println(l)
	}

	if err := t.SaveToFile(*flagPackingFile); err != nil {
		panic(fmt.Sprintf("%v", err))
	}
}
