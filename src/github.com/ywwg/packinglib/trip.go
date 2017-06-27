package packinglib

import (
	"fmt"
)

// Trip describes a trip, which includes a length and a context
type Trip struct {
	Days int
	C    *Context
}

type PackList map[string][]Item

var AllItems = make(PackList)

func RegisterItems(category string, items []Item) {
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

// MakeList returns a map of category to slice of PackedItems for the given trip
func (t *Trip) MakeList() PackList {
	packlist := make(PackList)
	for category, items := range AllItems {
		var packed []Item
		for _, i := range items {
			p := i.Pack(t)
			if p.Count() > 0 {
				packed = append(packed, p)
			}
		}
		packlist[category] = packed
	}

	return packlist
}
