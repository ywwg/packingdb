package packinglib

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

// PackList is a map from category name to slice of items
type PackList map[Category][]*Item

// XXXXXXXXXXXXXXXX OK the main thing we want to do is make this more of an API
// -- a trip is a thing we modify and then something about can ask questions
// about it or make mutations.
//
// Trip describes a trip, which includes a length and a context
type Trip struct {
	// list??? XXXX this happens because we build a big context out of little
	// ones... but basically what we should do is give the trip a context and a
	// property list I guess????  right now a trip is just a big context with a
	// bunch of bullshit convenience functions and a list of what is actually
	// packed (which itself is a convenience list).
	C           *Context
	contextName string

	// packList is a map of all the items in the trip.
	packList PackList
	// codeToItem is a map from a string code to the Item it corresponds to
	codeToItem map[string]*Item
	// itemToCode is the reverse.
	itemToCode map[*Item]string

	registry Registry
}

func getCode(idx int) string {
	code := ""
	adjust := 0
	for {
		codeVal := int32(idx%26 - adjust)
		code = string(rune('a')+codeVal) + code
		idx -= idx % 26
		idx /= 26
		if idx == 0 {
			break
		}
		adjust = 1
	}
	return code
}

// NewTripFromCustomContext returns a constructed trip for the given
// constructed context and number of nights.
func NewTripFromCustomContext(registry Registry, context *Context) (*Trip, error) {
	t := &Trip{
		C:           context,
		contextName: context.Name,
		registry:    registry,
	}
	t.packList = t.makeList()
	t.codeToItem = make(map[string]*Item)
	t.itemToCode = make(map[*Item]string)
	idx := 0
	keys := t.SortedCategories()

	for _, category := range keys {
		for _, item := range t.packList[category] {
			t.codeToItem[getCode(idx)] = item
			t.itemToCode[item] = getCode(idx)
			idx++
		}
	}
	return t, nil
}

// NewTrip returns a constructed trip for the given named context and
// number of nights.
func NewTrip(registry Registry, nights int, cname string) (*Trip, error) {
	c, err := registry.GetContext(cname)
	if err != nil {
		return nil, err
	}
	c.Nights = nights
	return NewTripFromCustomContext(registry, c)
}

func (t *Trip) AddProperty(p string) error {
	if err := t.C.addProperty(p); err != nil {
		return err
	}
	t.updateList()
	return nil
}

func (t *Trip) RemoveProperty(p string) error {
	if err := t.C.removeProperty(p); err != nil {
		return err
	}
	t.updateList()
	return nil
}

func (t *Trip) HasProperty(p Property) bool {
	return t.C.hasProperty(p)
}

// GetItemByCode returns the item with the given code, or false if not found.
func (t *Trip) GetItemByCode(code string) (*Item, bool) {
	item, ok := t.codeToItem[code]
	return item, ok
}

// makeList returns a map of category to slice of PackedItems for the given trip
func (t *Trip) makeList() PackList {
	packlist := make(PackList)
	for category, items := range t.registry.AllItems() {
		var toPack []*Item
		for _, i := range items {
			calced := i
			calced.AdjustCount(t.C)
			toPack = append(toPack, calced)
		}
		packlist[category] = toPack
	}

	return packlist
}

// updateList reAdjustCounts the packing list, changing the Count for each item.
// This can cause some items to go to Count of zero, but retain their packed
// state.  This allows users to remove and readd properties in one session
// and the packing state will be maintained. Once the list is written to disk
// however, the Count=0 items will be removed.
func (t *Trip) updateList() {
	for _, items := range t.packList {
		for _, i := range items {
			i.AdjustCount(t.C)
		}
	}
}

// Pack tries to pack the provided item first by short code, then by full name.
func (t *Trip) Pack(i string, pack bool) {
	// First try to pack by code
	if item, ok := t.codeToItem[i]; ok {
		item.Pack(pack)
		return
	}

	// Now fall back to string matching (which we do when loading the csv)
	found := false
	for _, items := range t.packList {
		for _, item := range items {
			if strings.EqualFold(item.Name(), i) {
				item.Pack(pack)
				found = true
			}
		}
	}
	if !found {
		fmt.Printf("tried to pack nonexistant item, ignoring: %s\n", i)
	}
}

// PackCategory packs the entire category.
func (t *Trip) PackCategory(cat string) {
	found := false
	for category := range t.packList {
		if strings.EqualFold(cat, string(category)) {
			found = true
			for _, i := range t.packList[category] {
				i.Pack(true)
			}
			break
		}
	}

	if !found {
		panic(fmt.Sprintf("tried to pack nonexistant category: %s", cat))
	}
}

// SortedCategories returns a list of all the categories currently represented
// by the context.
func (t *Trip) SortedCategories() []Category {
	// map iteration is nondeterministic so sort the keys.
	var keys []Category
	for k := range t.packList {
		keys = append(keys, k)
	}
	SortCategories(keys)
	return keys
}

// itemsNonEmpty returns true if at least one item in the list has a non-zero
// Count.
func itemsNonEmpty(items []*Item) bool {
	for _, i := range items {
		if i.Count() > 0 {
			return true
		}
	}
	return false
}

// Strings returns a slice of pretty strings representing the entire packing list.
func (t *Trip) Strings(showCat string, hideUnpacked bool) []string {
	var lines []string

	keys := t.SortedCategories()

	foundCat := false
	for _, category := range keys {
		if showCat != "" {
			if strings.EqualFold(string(category), showCat) {
				continue
			}
			foundCat = true
		}
		if itemsNonEmpty(t.packList[category]) {
			lines = append(lines, fmt.Sprintf("%s:", category))
		}
		for _, i := range t.packList[category] {
			if hideUnpacked && i.Packed() {
				continue
			}
			if i.Count() > 0 {
				lines = append(lines, fmt.Sprintf("\t%s", styleItem(i)))
			}
		}
	}
	if showCat != "" && !foundCat {
		panic(fmt.Sprintf("Didn't find category %s", showCat))
	}
	return lines
}

// PackingMenuItems returns a list of PackMenuItems for the given trip. Any
// categories in hiddenCategories will be hidden, and hidePacked will hide all
// packed items.
func (t *Trip) PackingMenuItems(hiddenCategories map[Category]bool, hidePacked bool) []PackMenuItem {
	var items []PackMenuItem
	keys := t.SortedCategories()

	for _, category := range keys {
		hide := hiddenCategories[category]
		if itemsNonEmpty(t.packList[category]) {
			var displayCat string
			if hide {
				displayCat = fmt.Sprintf("⊞ %s", category)
			} else {
				displayCat = fmt.Sprintf("⊟ %s", category)
			}
			items = append(items, NewMenuItem(displayCat, MenuCategory, string(category)))
		}
		if !hide {
			for _, i := range t.packList[category] {
				if hidePacked && i.Packed() {
					continue
				}
				if i.Count() > 0 {
					items = append(items, NewMenuItem(styleItem(i), MenuPackable, t.itemToCode[i]))
				}
			}
		}
	}
	return items
}

func styleItem(i *Item) string {
	checkbox := "○"
	if i.Packed() {
		checkbox = "●"
	}
	return fmt.Sprintf("%s %s", checkbox, i.String())
}

func (t *Trip) styleProperty(p Property) string {
	if t.HasProperty(p) {
		return fmt.Sprintf("● %-20s %s", string(p), t.registry.GetDescription(p))
	}
	return fmt.Sprintf("○ %-20s %s", string(p), t.registry.GetDescription(p))
}

// PropertyMenuItems returns a list of PackMenuItems for the given trip.
// Any categories in hiddenCategories will be hidden, and hidePacked will hide
// all packed items.
func (t *Trip) PropertyMenuItems() []PackMenuItem {
	var l []Property
	for name := range t.registry.AllProperties() {
		l = append(l, name)
	}
	less := func(i, j int) bool {
		return strings.ToLower(string(l[i])) < strings.ToLower(string(l[j]))
	}
	sort.Slice(l, less)

	var items []PackMenuItem
	for _, prop := range t.registry.ListProperties() {
		displayProp := t.styleProperty(prop)
		items = append(items, NewMenuItem(displayProp, MenuProperty, string(prop)))
	}
	return items
}

// ToggleItemPacked flips the packed state of the given item.
func (t *Trip) ToggleItemPacked(code string) error {
	// Only works with codes
	item, ok := t.codeToItem[code]
	if !ok {
		return fmt.Errorf("couldn't find item to pack with code %s", code)
	}

	item.Pack(!item.Packed())
	return nil
}

// LoadFromFile initializes the trip from the given file.
// Old file format:
// first line: number of nights, context name
// following lines: true/false string, name of packed item
//
// New file format:
// first line: "V2", number of nights, tmin, tmax, context name, contexts...
// if context_name is known, other contexts are added to it.
func LoadFromFile(r Registry, nights int, f string) (*Trip, error) {
	switch filepath.Ext(f) {
	case ".csv":
		return LoadFromCSV(r, nights, f)
	case ".yml", ".yaml":
		return LoadFromYAML(r, nights, f)
	default:
		return nil, fmt.Errorf("load extension not recognized: %s", filepath.Ext(f))
	}
}

func LoadFromCSV(r Registry, nights int, f string) (*Trip, error) {
	dat, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}
	t := &Trip{
		registry: r,
	}
	buf := bytes.NewBuffer(dat)
	scanner := bufio.NewScanner(buf)
	for i := 0; scanner.Scan(); i++ {
		if i == 0 {
			toks := strings.Split(scanner.Text(), ",")
			// V3 needs to store context names explicitly different from property names.
			if toks[0] == "V2" {
				var err error
				fileNights, err := strconv.Atoi(toks[1])
				if err != nil {
					return nil, err
				}
				if nights == 0 {
					nights = fileNights
				}
				tmin, err := strconv.Atoi(toks[2])
				if err != nil {
					return nil, err
				}
				tmax, err := strconv.Atoi(toks[3])
				if err != nil {
					return nil, err
				}
				context, err := t.registry.GetConcreteContext(toks[4], nights, tmin, tmax)
				if err != nil {
					context, err = NewContext(t.registry, toks[4], nights, tmin, tmax, nil)
				}
				if err != nil {
					panic(fmt.Sprintf("Error while building context for trip: %s", err.Error()))
				}
				for _, prop := range toks[5:] {
					if err := context.addProperty(prop); err != nil {
						panic(fmt.Sprintf("Error adding property while building trip: %s", err.Error()))
					}
				}
				loaded, err := NewTripFromCustomContext(t.registry, context)
				if err != nil {
					panic(err.Error())
				}
				*t = *loaded
			} else {
				if len(toks) != 2 {
					panic("Expected exactly two values in non-custom (old style) file")
				}
				nights, err := strconv.Atoi(toks[0])
				if err != nil {
					return nil, err
				}
				loaded, err := NewTrip(t.registry, nights, toks[1])
				if err != nil {
					panic(err.Error())
				}
				*t = *loaded
			}
		} else {
			toks := strings.SplitN(scanner.Text(), ",", 2)
			if toks[0] == "true" {
				t.Pack(toks[1], true)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("reading file: %s", err))
	}
	return t, nil
}

func LoadFromYAML(r Registry, nights int, f string) (*Trip, error) {
	reader, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	dec := yaml.NewDecoder(reader)
	var yt YamlTrip
	err = dec.Decode(&yt)
	if err != nil {
		return nil, err
	}
	if yt.Version != yaml_version {
		return nil, fmt.Errorf("incorrect version in yaml, want %d, got %d", yaml_version, yt.Version)
	}
	return yt.AsTrip(r)
}

// SaveToFile saves the trip to the provided filename.
func (t *Trip) SaveToFile(f string) error {
	switch filepath.Ext(f) {
	case ".csv":
		err := t.SaveToCSV(f)
		if err != nil {
			return err
		}
		return t.SaveToYAML(f + ".yml")
		// return t.SaveToCSV(f)
	case ".yml", ".yaml":
		return t.SaveToYAML(f)
	default:
		return fmt.Errorf("save extension not recognized")
	}
}

func (t *Trip) SaveToCSV(f string) error {
	packedcsv := fmt.Sprintf("V2,%d,%d,%d,%s,", t.C.Nights, t.C.TemperatureMin, t.C.TemperatureMax, t.contextName)
	for p, val := range t.C.Properties {
		if val {
			packedcsv += fmt.Sprintf("%s,", string(p))
		}
	}
	packedcsv += "\n"
	for _, items := range t.packList {
		for _, item := range items {
			if item.Count() > 0 {
				packedcsv += fmt.Sprintf("%v,%s\n", item.Packed(), item.Name())
			}
		}
	}
	return os.WriteFile(f, []byte(packedcsv), 0644)
}

func (t *Trip) SaveToYAML(f string) error {
	var buf []byte
	writer := bytes.NewBuffer(buf)
	yt := FromTrip(t)

	enc := yaml.NewEncoder(writer)
	err := enc.Encode(yt)
	if err != nil {
		return err
	}

	return os.WriteFile(f, writer.Bytes(), 0644)
}
