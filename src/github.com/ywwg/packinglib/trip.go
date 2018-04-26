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

// PackList is a map from category name to slice of items
type PackList map[string][]Item

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
	}
	if existing, ok := AllItems[category]; ok {
		AllItems[category] = append(existing, items...)
		return
	}
	AllItems[category] = items
}

var contexts = make(map[string]Context)

// RegisterContext registers the given context with the system.
// Also registers a property with the context name.
func RegisterContext(c Context) {
	if _, ok := contexts[c.Name]; ok {
		panic(fmt.Sprintf("Duplicate context: %s", c.Name))
	}
	contexts[c.Name] = c
	RegisterProperty(Property(c.Name))
}

func GetContext(name string) *Context {
	c := &Context{}
	found, ok := contexts[name]
	if !ok {
		panic(fmt.Sprintf("Unknown context: %s", name))
	}
	*c = found
	return c
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

func NewTrip(nights int, context string) *Trip {
	t := &Trip{
		Nights:      nights,
		C:           GetContext(context),
		contextName: context,
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
		for _, item := range t.packList[category] {
			t.codeToItem[getCode(idx)] = item
			t.itemToCode[item] = getCode(idx)
			idx++
		}
	}
	return t
}

// MakeList returns a map of category to slice of PackedItems for the given trip
func (t *Trip) makeList() PackList {
	packlist := make(PackList)
	for category, items := range AllItems {
		var toPack []Item
		for _, i := range items {
			calced := i.Itemize(t)
			if calced.Count() > 0 {
				toPack = append(toPack, calced)
			}
		}
		packlist[category] = toPack
	}

	return packlist
}

func (t *Trip) Pack(i string) {
	// First try to pack by code
	if item, ok := t.codeToItem[i]; ok {
		item.Pack()
		return
	}

	// Now fall back to string matching (which we do when loading the csv)
	found := false
	for _, items := range t.packList {
		for _, item := range items {
			if strings.ToLower(item.Name()) == strings.ToLower(i) {
				item.Pack()
				found = true
			}
		}
	}
	if !found {
		panic(fmt.Sprintf("tried to pack nonexistant item: %s", i))
	}
}

func (t *Trip) PackCategory(cat string) {
	if items, ok := t.packList[cat]; ok {
		for _, i := range items {
			i.Pack()
		}
	} else {
		panic(fmt.Sprintf("tried to pack nonexistant category: %s", cat))
	}
}

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
		if len(t.packList[category]) > 0 {
			lines = append(lines, fmt.Sprintf("%s:", category))
		}
		for _, i := range t.packList[category] {
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

func (t *Trip) LoadFromFile(f string) error {
	dat, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(dat)
	scanner := bufio.NewScanner(buf)
	for i := 0; scanner.Scan(); i++ {
		toks := strings.SplitN(scanner.Text(), ",", 2)
		if i == 0 {
			nights, err := strconv.Atoi(toks[0])
			if err != nil {
				return err
			}
			*t = *NewTrip(nights, toks[1])
		} else {
			if toks[0] == "true" {
				t.Pack(toks[1])
			}
		}
	}
	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("reading file: %s", err))
	}
	return nil
}

func (t *Trip) SaveToFile(f string) error {
	packedcsv := fmt.Sprintf("%d,%s\n", t.Nights, t.contextName)
	for _, items := range t.packList {
		for _, item := range items {
			packedcsv += fmt.Sprintf("%v,%s\n", item.Packed(), item.Name())
		}
	}
	return ioutil.WriteFile(f, []byte(packedcsv), 0644)
}
