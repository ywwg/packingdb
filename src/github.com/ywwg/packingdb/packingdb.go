// packingdb is the newer interactive readline-based application for
// creating and editing packing lists.  It is currently the recommended binary.
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/ywwg/packinglib"
	"golang.org/x/crypto/ssh/terminal"

	_ "github.com/ywwg/contexts"
	_ "github.com/ywwg/items"
)

func mainMenu(t *packinglib.Trip) (string, error) {
	items := []string{
		"Pack",
		"Configure Trip Properties",
		"Set Nights",
		"Set Min Temperature",
		"Set Max Temperature",
		"Quit",
	}

	prompt := promptui.Select{
		Label: "Main Menu",
		Items: items,
		Size:  len(items),
	}

	_, result, err := prompt.Run()

	return result, err
}

func nameEntry() (string, error) {
	validate := func(input string) error {
		if len(input) == 0 {
			return errors.New("Must provide a name for this trip")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    "Trip Name",
		Validate: validate,
	}

	return prompt.Run()
}

func nightsEntry(current int) (int, error) {
	validate := func(input string) error {
		if input == "" {
			return nil
		}
		_, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return errors.New("Invalid number")
		}
		return nil
	}
	label := "Number of nights"
	if current != -1 {
		label = fmt.Sprintf("Number of nights (current setting: %d)", current)
	}
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return -1, err
	}
	if result == "" {
		return current, nil
	}
	temp, err := strconv.ParseInt(result, 10, 64)
	return int(temp), err
}

func temperatureEntry(help string, current int) (int, error) {
	validate := func(input string) error {
		if input == "" {
			return nil
		}
		_, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return errors.New("Invalid number")
		}
		return nil
	}
	label := help + " Temperature"
	if current != -100 {
		label = fmt.Sprintf("%s Temperature (current setting: %d)", help, current)
	}
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return -1, err
	}
	if result == "" {
		return current, nil
	}
	temp, err := strconv.ParseInt(result, 10, 64)
	return int(temp), err
}

func packMenu(t *packinglib.Trip) error {
	cursor := 0
	hidePacked := false
	hiddenCats := make(map[string]bool)

	BackMenuItem := packinglib.NewMenuItem("↩ Back", packinglib.MenuAction, "back")
	HidePackedMenuItem := packinglib.NewMenuItem("Hide Packed", packinglib.MenuAction, "hidepacked")
	ShowPackedMenuItem := packinglib.NewMenuItem("Show Packed", packinglib.MenuAction, "hidepacked")
	UnhideAllCatsItem := packinglib.NewMenuItem("Show All Categories", packinglib.MenuAction, "showcats")

	for {
		items := []packinglib.PackMenuItem{BackMenuItem}
		items = append(items, UnhideAllCatsItem)
		if hidePacked {
			items = append(items, ShowPackedMenuItem)
		} else {
			items = append(items, HidePackedMenuItem)
		}
		items = append(items, t.MenuItems(hiddenCats, hidePacked)...)
		_, height, err := terminal.GetSize(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatal(fmt.Sprintf("couldn't get terminal size %v", err))
		}
		// Take off 10 lines to account for bits of ui (history, etc)
		if height > 10 {
			height -= 10
		}
		prompt := promptui.Select{
			Label: "Packing Menu",
			Items: items,
			Templates: &promptui.SelectTemplates{
				Label:    "  {{ .Name }}",
				Active:   "▸ {{ .Name | underline}}",
				Inactive: "  {{ .Name }}",
				Selected: "{{ .Name }}",
			},
			Size: height,
			// TODO: figure out how to bind to pageup and pagedown
			Keys: &promptui.SelectKeys{
				Prev:     promptui.Key{Code: promptui.KeyPrev, Display: promptui.KeyPrevDisplay},
				Next:     promptui.Key{Code: promptui.KeyNext, Display: promptui.KeyNextDisplay},
				PageUp:   promptui.Key{Code: promptui.KeyBackward, Display: promptui.KeyBackwardDisplay},
				PageDown: promptui.Key{Code: promptui.KeyForward, Display: promptui.KeyForwardDisplay},
				Search:   promptui.Key{Code: '/', Display: "/"},
			},
		}

		i, _, err := prompt.RunStartingAt(cursor)
		if err != nil {
			return err
		}
		selected := items[i]
		// TODO: I forget, does go have operator overriding? Probs not.
		if selected.Equals(BackMenuItem) {
			return nil
		} else if selected.Equals(HidePackedMenuItem) {
			hidePacked = !hidePacked
		} else if selected.Equals(UnhideAllCatsItem) {
			hiddenCats = make(map[string]bool)
		} else if selected.Type == packinglib.MenuCategory {
			hiddenCats[selected.Code] = !hiddenCats[selected.Code]
		} else if selected.Type == packinglib.MenuPackable {
			if err := t.ToggleItemPacked(selected.Code); err != nil {
				log.Fatal(err)
			}
		}
		cursor = i
	}
}

func styleProperty(t *packinglib.Trip, prop string) string {
	if t.C.HasProperty(prop) {
		return "☑ " + prop
	}
	return "☐ " + prop
}

func propertyMenu(t *packinglib.Trip) error {
	cursor := 0

	for {
		items := []string{"↩ Back"}
		for _, prop := range packinglib.ListProperties() {
			items = append(items, styleProperty(t, prop))
		}
		_, height, err := terminal.GetSize(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatal(fmt.Sprintf("couldn't get terminal size %v", err))
		}
		// Take off 10 lines to account for bits of ui (history, etc)
		if height > 10 {
			height -= 10
		}
		prompt := promptui.Select{
			Label: "Property Menu",
			Items: items,
			Size:  height,
		}

		i, _, err := prompt.RunStartingAt(cursor)
		if err != nil {
			return err
		}

		// Hack to remove the checkbox before the name of the property
		selected := strings.SplitAfterN(items[i], " ", 2)[1]
		if selected == "Back" {
			return nil
		} else if t.C.HasProperty(selected) {
			if err := t.C.RemoveProperty(selected); err != nil {
				log.Fatal(err)
			}
		} else {
			if err := t.C.AddProperty(selected); err != nil {
				log.Fatal(err)
			}
		}
		cursor = i
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("Need a packing filename (one will be created if it doesn't exist)")
	}
	filename := args[0]

	// File mode: load if the file already exists or create new if not
	t := &packinglib.Trip{}
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		if err2 := t.LoadFromFile(0, filename); err2 != nil {
			log.Fatal(fmt.Sprintf("%v", err2))
		}
	} else {
		title, err := nameEntry()
		if err != nil {
			log.Fatal(err)
		}

		nights, err := nightsEntry(-1)
		if err != nil {
			log.Fatal(err)
		}

		context := &packinglib.Context{
			Name:       title,
			Properties: make(packinglib.PropertySet),
		}
		context.TemperatureMin, err = temperatureEntry("Minimum", -100)
		if err != nil {
			log.Fatal(err)
		}
		context.TemperatureMax, err = temperatureEntry("Maximum", -100)
		if err != nil {
			log.Fatal(err)
		}

		t, err = packinglib.NewTripFromCustomContext(nights, context)
		if err != nil {
			log.Fatal(err)
		}

		propertyMenu(t)
	}

	// Main Menu
	for {
		fmt.Printf("\033[1m")
		fmt.Printf("\n\nName: %s\n", t.C.Name)
		fmt.Printf("Nights: %d\n", t.Nights)
		fmt.Printf("Temperatures: %d - %d\n", t.C.TemperatureMin, t.C.TemperatureMax)
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
		fmt.Printf("\n")
		fmt.Printf("\033[0m")

		result, err := mainMenu(t)
		if err != nil {
			break
		}
		if result == "Quit" {
			break
		}
		switch result {
		case "Configure Trip Properties":
			propertyMenu(t)
		case "Set Nights":
			if nights, err := nightsEntry(t.Nights); err == nil {
				t.Nights = nights
			}
		case "Set Min Temperature":
			if temp, err := temperatureEntry("Minimum", t.C.TemperatureMin); err == nil {
				t.C.TemperatureMin = temp
			}
		case "Set Max Temperature":
			if temp, err := temperatureEntry("Maximum", t.C.TemperatureMax); err == nil {
				t.C.TemperatureMax = temp
			}
		case "Pack":
			err := packMenu(t)
			if err != nil {
				break
			}
		}
	}

	if err := t.SaveToFile(filename); err != nil {
		log.Fatal(fmt.Sprintf("%v", err))
	}
}
