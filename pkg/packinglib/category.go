package packinglib

import (
	"sort"
	"strings"
)

// Category is a name for what kind of thing an item is being packed.
type Category string

func SortCategories(cats []Category) {
	less := func(i, j int) bool {
		return strings.ToLower(string(cats[i])) < strings.ToLower(string(cats[j]))
	}
	sort.Slice(cats, less)
}
