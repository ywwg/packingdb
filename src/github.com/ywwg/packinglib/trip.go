package packinglib

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

type PackCategory struct {
	Visible bool
	Items   []Item
}

// PackList is a map from category name to slice of items
type PackList map[string]*PackCategory

func (p PackList) ItemsForCategory(category string) []Item {
	return p[category].Items
}

// AllItems is a convenience map of all items that packingdb knows about
var AllItems = make(PackList)

// Trip describes a trip, which includes a length and a context
type Trip struct {
	Nights      int
	C           *Context
	contextName string
	// packList is a map of all the items in the trip.
	packList PackList
	// codeToItem is a map from a string code to the Item it corresponds to
	codeToItem map[string]Item
	// itemToCode is the reverse.
	itemToCode map[Item]string
}

// dupeChecker is a map to track all of the item names and make sure we don't
// have any duplicates.
var dupeChecker = make(map[string]bool)

// RegisterItems appends the given slice of Items to the registry under
// the given category.  Duplicate categories will be appended.  Items
// with duplicate case-insensitive names, even across categories,
// cause a panic.
func RegisterItems(category string, items []Item) {
	for _, i := range items {
		if _, ok := dupeChecker[strings.ToLower(i.Name())]; ok {
			panic(fmt.Sprintf("Duplicate item name: %s: %s", category, i.Name()))
		}
		dupeChecker[i.Name()] = true
		for p := range i.Prerequisites() {
			if _, ok := allProperties[p]; !ok {
				panic(fmt.Sprintf("Prerequisite property not found in allProperties, is it registered?: %s", p))
			}
		}
	}
	if existing, ok := AllItems[category]; ok {
		existing.Items = append(existing.Items, items...)
		return
	}
	AllItems[category] = &PackCategory{
		Visible: true,
		Items:   items,
	}
}

// Context is struct that holds data about the context of the trip
type Context struct {
	// Name of the context ("The Cape", "The Tiny House", "Firefly")
	Name string

	// TemperatureMin is the anticipated minimum temperature.
	TemperatureMin int

	// TemperatureMax is the anticipated maximum temperature.
	TemperatureMax int

	Properties PropertySet
}

var contexts = make(map[string]Context)

// ContextList returns a sorted slice of strings of the contexts.
func ContextList() []string {
	keys := make([]string, len(contexts))
	i := 0
	for k := range contexts {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

// RegisterContext registers the given context with the system.
// Also registers a property with the context name.
func RegisterContext(c Context) {
	if _, ok := contexts[c.Name]; ok {
		panic(fmt.Sprintf("Duplicate context: %s", c.Name))
	}
	contexts[c.Name] = c
	RegisterProperty(Property(c.Name))
}

// GetContext returns the context of the given name, or panics if not found.
func GetContext(name string) (*Context, error) {
	c := &Context{}
	found, ok := contexts[name]
	if !ok {
		return nil, fmt.Errorf("Unknown context: %s", name)
	}
	*c = found
	return c, nil
}

// GetContextWithTemperature loads the given context and substitutes the provided
// temperature range.
func GetContextWithTemperature(name string, tmin, tmax int) (*Context, error) {
	c, err := GetContext(name)
	if err != nil {
		return nil, err
	}
	c.TemperatureMin = tmin
	c.TemperatureMax = tmax
	return c, nil
}

// NewContext creates a new context with the given name, temperature range, and properties.
// Returns nil if any of the properties is unknown.  Properties are optional.
func NewContext(name string, tmin, tmax int, properties []string) (*Context, error) {
	c := &Context{
		Name:           name,
		TemperatureMin: tmin,
		TemperatureMax: tmax,
		Properties:     make(PropertySet),
	}

	RegisterProperty(Property(c.Name))
	if err := c.AddProperty(c.Name); err != nil {
		return nil, err
	}
	for _, p := range properties {
		if err := c.AddProperty(p); err != nil {
			return nil, err
		}
	}

	RegisterContext(*c)
	return c, nil
}

// AddProperty adds the property with the given name to the context, or returns
// error if it's not found. Empty strings are ignored.
func (c *Context) AddProperty(prop string) error {
	if prop == "" {
		return nil
	}
	if _, ok := allProperties[Property(prop)]; !ok {
		return fmt.Errorf("didn't find property, is it registered?: %s", prop)
	}
	c.Properties[Property(prop)] = true
	return nil
}

func getCode(idx int) string {
	code := ""
	adjust := 0
	for {
		codeVal := idx%26 - adjust
		code = string('a'+codeVal) + code
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
func NewTripFromCustomContext(nights int, context *Context) (*Trip, error) {
	t := &Trip{
		Nights:      nights,
		C:           context,
		contextName: context.Name,
	}
	t.packList = t.makeList()
	t.codeToItem = make(map[string]Item)
	t.itemToCode = make(map[Item]string)
	idx := 0
	var keys []string
	for k := range t.packList {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, category := range keys {
		for _, item := range t.packList.ItemsForCategory(category) {
			t.codeToItem[getCode(idx)] = item
			t.itemToCode[item] = getCode(idx)
			idx++
		}
	}
	return t, nil
}

// NewTrip returns a constructed trip for the given named context and
// number of nights.
func NewTrip(nights int, cname string) (*Trip, error) {
	c, err := GetContext(cname)
	if err != nil {
		return nil, err
	}
	return NewTripFromCustomContext(nights, c)
}

// MakeList returns a map of category to slice of PackedItems for the given trip
func (t *Trip) makeList() PackList {
	packlist := make(PackList)
	for category, packCat := range AllItems {
		var toPack []Item
		for _, i := range packCat.Items {
			calced := i.Itemize(t)
			if calced.Count() > 0 {
				toPack = append(toPack, calced)
			}
		}
		packlist[category] = &PackCategory{
			Visible: true,
			Items:   toPack,
		}
	}

	return packlist
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
	for _, packCat := range t.packList {
		for _, item := range packCat.Items {
			if strings.ToLower(item.Name()) == strings.ToLower(i) {
				item.Pack(pack)
				found = true
			}
		}
	}
	if !found {
		panic(fmt.Sprintf("tried to pack nonexistant item: %s", i))
	}
}

// PackCategory packs the entire category.
func (t *Trip) PackCategory(cat string) {
	found := false
	for category := range t.packList {
		if strings.ToLower(cat) == strings.ToLower(category) {
			found = true
			for _, i := range t.packList.ItemsForCategory(category) {
				i.Pack(true)
			}
			break
		}
	}

	if !found {
		panic(fmt.Sprintf("tried to pack nonexistant category: %s", cat))
	}
}

// Strings returns a slice of pretty strings representing the entire packing list.
func (t *Trip) Strings(showCat string, hideUnpacked bool) []string {
	var lines []string
	// map iteration is nondeterministic so sort the keys.
	var keys []string
	for k := range t.packList {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	foundCat := false
	for _, category := range keys {
		if showCat != "" {
			if strings.ToLower(category) != strings.ToLower(showCat) {
				continue
			}
			foundCat = true
		}
		if len(t.packList.ItemsForCategory(category)) > 0 {
			lines = append(lines, fmt.Sprintf("%s:", category))
		}
		for _, i := range t.packList.ItemsForCategory(category) {
			if hideUnpacked && i.Packed() {
				continue
			}
			lines = append(lines, fmt.Sprintf("\t(%s) %s", t.itemToCode[i], i.String()))
		}
	}
	if showCat != "" && !foundCat {
		panic(fmt.Sprintf("Didn't find category %s", showCat))
	}
	return lines
}

func (t *Trip) MenuItems() []PackMenuItem {
	var items []PackMenuItem
	// map iteration is nondeterministic so sort the keys.
	var keys []string
	for k := range t.packList {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, category := range keys {
		if len(t.packList.ItemsForCategory(category)) > 0 {
			items = append(items, New(category, MenuCategory, category))
		}
		if t.packList[category].Visible {
			for _, i := range t.packList.ItemsForCategory(category) {
				items = append(items, New(i.String(), MenuPackable, t.itemToCode[i]))
			}
		}
	}
	return items
}

func (t *Trip) ToggleCategoryVisibility(cat string) error {
	pl, ok := t.packList[cat]
	if !ok {
		return fmt.Errorf("didn't find category %s", cat)
	}
	t.packList[cat].Visible = !pl.Visible
	return nil
}

func (t *Trip) ToggleItemPacked(code string) error {
	// Only works with codes
	item, ok := t.codeToItem[code]
	if !ok {
		return fmt.Errorf("Couldn't find item to pack with code %s", code)
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
func (t *Trip) LoadFromFile(nights int, f string) error {
	dat, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(dat)
	scanner := bufio.NewScanner(buf)
	for i := 0; scanner.Scan(); i++ {
		if i == 0 {
			toks := strings.Split(scanner.Text(), ",")
			if toks[0] == "V2" {
				var err error
				fileNights, err := strconv.Atoi(toks[1])
				if err != nil {
					return err
				}
				if nights == 0 {
					nights = fileNights
				}
				tmin, err := strconv.Atoi(toks[2])
				if err != nil {
					return err
				}
				tmax, err := strconv.Atoi(toks[3])
				if err != nil {
					return err
				}
				var context *Context
				if _, ok := contexts[toks[4]]; ok {
					context, err = GetContextWithTemperature(toks[4], tmin, tmax)
				} else {
					context, err = NewContext(toks[4], tmin, tmax, nil)
				}
				if err != nil {
					panic(fmt.Sprintf("Error while building context for trip: %s", err.Error()))
				}
				for _, prop := range toks[5:] {
					if err := context.AddProperty(prop); err != nil {
						panic(fmt.Sprintf("Error adding property while building trip: %s", err.Error()))
					}
				}
				loaded, err := NewTripFromCustomContext(nights, context)
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
					return err
				}
				loaded, err := NewTrip(nights, toks[1])
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
	return nil
}

// SaveToFile saves the trip to the provided filename.
func (t *Trip) SaveToFile(f string) error {
	packedcsv := fmt.Sprintf("V2,%d,%d,%d,%s,", t.Nights, t.C.TemperatureMin, t.C.TemperatureMax, t.contextName)
	for p := range t.C.Properties {
		packedcsv += fmt.Sprintf("%s,", string(p))
	}
	packedcsv += "\n"
	for _, items := range t.packList {
		for _, item := range items.Items {
			packedcsv += fmt.Sprintf("%v,%s\n", item.Packed(), item.Name())
		}
	}
	return ioutil.WriteFile(f, []byte(packedcsv), 0644)
}
