package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/ywwg/packinglib"

	_ "github.com/ywwg/contexts"
	_ "github.com/ywwg/items"
)

var (
	flagContext        = flag.String("context", "", "The context you want to load or the name for a new empty context (not needed if packing file exists)")
	flagNights         = flag.Int("nights", 0, "The number of nights for the trip (not needed if packing file exists)")
	flagPackingFile    = flag.String("packfile", "", "The filename to create or load (not needed if you just want to print a list)")
	flagCategory       = flag.String("category", "", "Only print out the given category.")
	flagPackItem       = flag.String("pack", "", "The code of an item that has been packed (noop without packfile)")
	flagPackCategory   = flag.String("pack_category", "", "The name of an entire category that has been packed (noop without packfile)")
	flagHidePacked     = flag.Bool("hide_packed", false, "Only show unpacked items")
	flagListContexts   = flag.Bool("list_contexts", false, "List the available contexts and exit")
	flagListProperties = flag.Bool("list_properties", false, "List the available properties and exit")
	flagTempMin        = flag.Int("temp_min", math.MinInt64, "Set min temperature bound")
	flagTempMax        = flag.Int("temp_max", math.MaxInt64, "Set max temperature bound")
	flagAddProperty    = flag.String("add_property", "", "Add this property to the context")
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
			fmt.Printf("Context %s not found, creating empty context\n", *flagContext)
			c, err := packinglib.NewContext(*flagContext, -120, 120, nil)
			if err != nil {
				panic(err.Error())
			}
			t, err = packinglib.NewTripFromCustomContext(*flagNights, c)
			if err != nil {
				panic(err.Error())
			}
		}
	}

	if *flagTempMin != math.MinInt64 {
		t.C.TemperatureMin = *flagTempMin
	}
	if *flagTempMax != math.MaxInt64 {
		t.C.TemperatureMax = *flagTempMax
	}
	if *flagAddProperty != "" {
		for _, prop := range strings.Split(*flagAddProperty, ",") {
			if err := t.C.AddProperty(prop); err != nil {
				panic(err.Error())
			}
		}
	}

	if *flagPackItem != "" {
		t.Pack(*flagPackItem)
	}
	if *flagPackCategory != "" {
		t.PackCategory(*flagPackCategory)
	}

	// Print out trip context and information
	fmt.Printf("Context: %s\n", t.C.Name)
	fmt.Printf("Nights: %d\n", t.Nights)
	fmt.Printf("Properties: ")
	var plist []string
	for p := range t.C.Properties {
		plist = append(plist, string(p))
	}
	sort.Strings(plist)
	for i, p := range plist {
		if i == 0 {
			fmt.Printf("%s", p)
		} else {
			fmt.Printf(", %s", p)
		}
	}
	fmt.Printf("\n\n")

	// Print the packing list
	for _, l := range t.Strings(*flagCategory, *flagHidePacked) {
		fmt.Println(l)
	}

	if err := t.SaveToFile(*flagPackingFile); err != nil {
		panic(fmt.Sprintf("%v", err))
	}
}
