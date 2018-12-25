package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/ywwg/packinglib"
	"golang.org/x/crypto/ssh/terminal"

	_ "github.com/ywwg/contexts"
	_ "github.com/ywwg/items"
)

var (
	flagPackingFile = flag.String("packfile", "", "The filename to create or load (not needed if you just want to print a list)")
)

func main() {
	flag.Parse()

	if *flagPackingFile == "" {
		panic("Need a packing file")
	}

	// File mode: load if the file already exists or create new if not
	t := &packinglib.Trip{}
	if _, err := os.Stat(*flagPackingFile); os.IsNotExist(err) {
		panic("File not found")
	}

	if err := t.LoadFromFile(0, *flagPackingFile); err != nil {
		panic(fmt.Sprintf("%v", err))
	}

	// Main Menu
	for {
		result, err := mainMenu(t)
		if err != nil {
			break
		}
		if result == "Quit" {
			break
		}
		if result == "Pack" {
			err := packMenu(t)
			if err != nil {
				break
			}
		}
	}

	if err := t.SaveToFile(*flagPackingFile); err != nil {
		panic(fmt.Sprintf("%v", err))
	}
}

func mainMenu(t *packinglib.Trip) (string, error) {
	prompt := promptui.Select{
		Label: "Main Menu",
		Items: []string{
			"Pack",
			"Quit",
		},
	}

	_, result, err := prompt.Run()

	return result, err
}

func packMenu(t *packinglib.Trip) error {
	cursor := 0
	hidePacked := false
	hiddenCats := make(map[string]bool)

	BackMenuItem := packinglib.NewMenuItem("Back", packinglib.MenuAction, "back")
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
			panic(fmt.Sprintf("couldn't get terminal size %v", err))
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
				panic(err)
			}
		}
		cursor = i
	}
}
