package packinglib

import (
	//	"fmt"
	"sort"
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

// MakeList returns a map of category to slice of PackedItems for the given trip
func (t *Trip) MakeList() PackList {
	// map iteration is nondeterministic so sort the keys.
	var keys []string
	for k := range AllItems {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	packlist := make(PackList)
	for _, category := range keys {
		var packed []Item
		for _, i := range AllItems[category] {
			p := i.Pack(t)
			if p.Count() > 0 {
				packed = append(packed, p)
			}
		}
		packlist[category] = packed
	}

	return packlist
}
