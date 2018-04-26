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
	flagContext      = flag.String("context", "", "The context you want to load (not needed if packing file exists)")
	flagNights       = flag.Int("nights", 0, "The number of nights for the trip (not needed if packing file exists)")
	flagPackingFile  = flag.String("packfile", "", "The filename to create or load (not needed if you just want to print a list)")
	flagCategory     = flag.String("category", "", "Only print out the given category.")
	flagPackItem     = flag.String("pack", "", "The name of an item that has been packed (noop without packfile)")
	flagPackCategory = flag.String("pack_category", "", "The name of an entire category that has been packed (noop without packfile)")
	flagHidePacked   = flag.Bool("hide_packed", false, "Only show unpacked items")
)

func main() {
	flag.Parse()

	var t *packinglib.Trip

	// Simple mode: just print the list.
	if *flagPackingFile == "" {
		if len(*flagContext) == 0 {
			panic("Need a context")
		}
		if *flagNights == 0 {
			panic("Need a number of nights")
		}
		t = packinglib.NewTrip(*flagNights, *flagContext)
		for _, l := range t.Strings(*flagCategory, *flagHidePacked) {
			fmt.Println(l)
		}
		return
	}

	// File mode: load if the file already exists or create new if not
	if _, err := os.Stat(*flagPackingFile); !os.IsNotExist(err) {
		if len(*flagContext) != 0 {
			fmt.Println("(Ignoring context when loading file)")
		}
		if *flagNights != 0 {
			fmt.Println("(Ignoring nights when loading file)")
		}
		t = &packinglib.Trip{}
		if err2 := t.LoadFromFile(*flagPackingFile); err2 != nil {
			panic(fmt.Sprintf("%v", err2))
		}
	} else {
		if len(*flagContext) == 0 {
			panic("Need a context")
		}
		if *flagNights == 0 {
			panic("Need a number of nights")
		}
		t = packinglib.NewTrip(*flagNights, *flagContext)
	}

	if *flagPackItem != "" {
		t.Pack(*flagPackItem)
	}
	if *flagPackCategory != "" {
		t.PackCategory(*flagPackCategory)
	}
	for _, l := range t.Strings(*flagCategory, *flagHidePacked) {
		fmt.Println(l)
	}

	if err := t.SaveToFile(*flagPackingFile); err != nil {
		panic(fmt.Sprintf("%v", err))
	}
}
