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

// Trip describes a trip, which includes a length and a context
type Trip struct {
	Days        int
	C           *Context
	contextName string
	Packed      PackList
}

type PackList map[string][]Item

var AllItems = make(PackList)

// dupeChecker is a simple map to track all of the item names and
// make sure we don't have any duplicates.
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

func RegisterContext(name string, c Context) {
	if _, ok := contexts[name]; ok {
		panic(fmt.Sprintf("Duplicate context: %s", name))
	}
	contexts[name] = c
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

func NewTrip(days int, context string) *Trip {
	t := &Trip{
		Days:        days,
		C:           GetContext(context),
		contextName: context,
	}
	t.Packed = t.makeList()
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
	found := false
	for _, items := range t.Packed {
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

func (t *Trip) Strings(showCat string) []string {
	var lines []string
	// map iteration is nondeterministic so sort the keys.
	var keys []string
	for k := range t.Packed {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, category := range keys {
		if showCat != "" && strings.ToLower(category) != strings.ToLower(showCat) {
			continue
		}
		if len(t.Packed[category]) > 0 {
			lines = append(lines, fmt.Sprintf("%s:", category))
		}
		for _, i := range t.Packed[category] {
			lines = append(lines, fmt.Sprintf("\t%s", i.String()))
		}
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
			days, err := strconv.Atoi(toks[0])
			if err != nil {
				return err
			}
			*t = *NewTrip(days, toks[1])
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
	packedcsv := fmt.Sprintf("%d,%s\n", t.Days, t.contextName)
	for _, items := range t.Packed {
		for _, item := range items {
			packedcsv += fmt.Sprintf("%v,%s\n", item.Packed(), item.Name())
		}
	}
	return ioutil.WriteFile(f, []byte(packedcsv), 0644)
}
