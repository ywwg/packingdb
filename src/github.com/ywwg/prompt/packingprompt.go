package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/ywwg/packinglib"

	_ "github.com/ywwg/contexts"
	_ "github.com/ywwg/items"
)

var (
	flagContext        = flag.String("context", "", "The context you want to load or the name for a new empty context (not needed if packing file exists)")
	flagNights         = flag.Int("nights", 0, "The number of nights for the trip (not needed if packing file exists)")
	flagPackingFile    = flag.String("packfile", "", "The filename to create or load (not needed if you just want to print a list)")
	flagListContexts   = flag.Bool("list-contexts", false, "List the available contexts and exit")
	flagListProperties = flag.Bool("list-properties", false, "List the available properties and exit")
	flagMinTemp        = flag.Int("min-temp", math.MinInt64, "Set min temperature bound")
	flagMaxTemp        = flag.Int("max-temp", math.MaxInt64, "Set max temperature bound")
	flagAddProperty    = flag.String("add-property", "", "Add this property to the context")
)

func main() {
	flag.Parse()

	if *flagPackingFile == "" {
		panic("Need a packing file")
	}

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

	// File mode: load if the file already exists or create new if not
	t = &packinglib.Trip{}
	if _, err := os.Stat(*flagPackingFile); !os.IsNotExist(err) {
		if len(*flagContext) != 0 {
			fmt.Println("(Ignoring context when loading file)")
		}
		if err2 := t.LoadFromFile(*flagNights, *flagPackingFile); err2 != nil {
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

	if *flagMinTemp != math.MinInt64 {
		t.C.TemperatureMin = *flagMinTemp
	}
	if *flagMaxTemp != math.MaxInt64 {
		t.C.TemperatureMax = *flagMaxTemp
	}
	if *flagAddProperty != "" {
		for _, prop := range strings.Split(*flagAddProperty, ",") {
			if err := t.C.AddProperty(prop); err != nil {
				panic(err.Error())
			}
		}
	}

	// Main Menu
	for {
		result, err := mainMenu()
		if err != nil {
			break
		}
		if result == "Quit" {
			break
		}
		if result == "Pack" {
			err := packMenu()
			if err != nil {
				break
			}
		}
	}

	if err := t.SaveToFile(*flagPackingFile); err != nil {
		panic(fmt.Sprintf("%v", err))
	}
}

func mainMenu() (string, error) {
	prompt := promptui.Select{
		Label: "Main Menu",
		Items: []string{"Pack", "Quit"},
	}

	_, result, err := prompt.Run()

	return result, err
}

func packMenu() error {
	for {
		prompt := promptui.Select{
			Label: "Packing Menu",
			Items: []string{"Back", "thing"},
		}

		_, result, err := prompt.Run()
		if err != nil {
			return err
		}
		if result == "Back" {
			return nil
		}
	}
}
