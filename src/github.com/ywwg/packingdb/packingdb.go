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
	flagContext        = flag.String("context", "", "The context you want to load (not needed if packing file exists)")
	flagNights         = flag.Int("nights", 0, "The number of nights for the trip (not needed if packing file exists)")
	flagPackingFile    = flag.String("packfile", "", "The filename to create or load (not needed if you just want to print a list)")
	flagCategory       = flag.String("category", "", "Only print out the given category.")
	flagPackItem       = flag.String("pack", "", "The code of an item that has been packed (noop without packfile)")
	flagPackCategory   = flag.String("pack_category", "", "The name of an entire category that has been packed (noop without packfile)")
	flagHidePacked     = flag.Bool("hide_packed", false, "Only show unpacked items")
	flagListContexts   = flag.Bool("list_contexts", false, "List the available contexts and exit")
	flagListProperties = flag.Bool("list_properties", false, "List the available properties and exit")
)

func main() {
	flag.Parse()

	var t *packinglib.Trip
	if *flagListContexts {
		for _, c := range packinglib.ContextList() {
			fmt.Println(c)
		}
		return
	}
	if *flagListProperties {
		for _, p := range packinglib.ListProperties() {
			fmt.Println(p)
		}
		return
	}

	// Simple mode: just print the list.
	if *flagPackingFile == "" {
		if len(*flagContext) == 0 {
			panic("Need a context")
		}
		if *flagNights == 0 {
			panic("Need a number of nights")
		}
		t, err := packinglib.NewTrip(*flagNights, *flagContext)
		if err != nil {
			panic(err.Error())
		}
		for _, l := range t.Strings(*flagCategory, *flagHidePacked) {
			fmt.Println(l)
		}
		return
	}

	// File mode: load if the file already exists or create new if not
	t = &packinglib.Trip{}
	if _, err := os.Stat(*flagPackingFile); !os.IsNotExist(err) {
		if len(*flagContext) != 0 {
			fmt.Println("(Ignoring context when loading file)")
		}
		if *flagNights != 0 {
			fmt.Println("(Ignoring nights when loading file)")
		}
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
		var err error
		t, err = packinglib.NewTrip(*flagNights, *flagContext)
		if err != nil {
			panic(err.Error())
		}
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
