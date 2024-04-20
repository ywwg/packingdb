package packinglib

import (
	"fmt"
	"sort"
	"strings"
)

// ContextRegistry is a generic interface for uhhh doing things with contexts.
type Registry interface {
	ContextList() []string

	// Not clear how registration should work... probably just "add"
	RegisterContext(c Context)
	RegisterItems(category Category, items []*Item)
	GetContext(name string) (*Context, error)
	GetContextTemperatureRange(name string, tmin, tmax int) (*Context, error)
	HasProperty(p Property) bool
}

// StructRegistry is a simple in-memory impl of a ContextRegistry
type StructRegistry struct {
	contexts      map[string]Context
	allProperties map[Property]string
	AllItems      PackList
}

// ContextList returns a sorted slice of strings of the contexts.
func (t *StructRegistry) ContextList() []string {
	keys := make([]string, len(t.contexts))
	i := 0
	for k := range t.contexts {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

// RegisterContext registers the given context with the system.
// Also registers a property with the context name.
func (t *StructRegistry) RegisterContext(c Context) {
	if _, ok := t.contexts[c.Name]; ok {
		panic(fmt.Sprintf("Duplicate context: %s", c.Name))
	}
	t.contexts[c.Name] = c
	RegisterProperty(Property(c.Name), "")
}

// GetContext returns the context of the given name, or returns error if not found.
func (t *StructRegistry) GetContext(name string) (*Context, error) {
	c := &Context{}
	found, ok := t.contexts[name]
	if !ok {
		return nil, fmt.Errorf("unknown context: %s", name)
	}
	*c = found
	return c, nil
}

// GetContextTemperatureRange loads the given context and substitutes the provided
// temperature range.
func (t *StructRegistry) GetContextTemperatureRange(name string, tmin, tmax int) (*Context, error) {
	c, err := t.GetContext(name)
	if err != nil {
		return nil, err
	}
	c.TemperatureMin = tmin
	c.TemperatureMax = tmax
	return c, nil
}

// RegisterItems appends the given slice of Items to the registry under
// the given category.  Duplicate categories will be appended.  Items
// with duplicate case-insensitive names, even across categories,
// cause a panic.
func (t *StructRegistry) RegisterItems(category Category, items []*Item) {
	for _, i := range items {
		if _, ok := dupeChecker[strings.ToLower(i.Name())]; ok {
			panic(fmt.Sprintf("Duplicate item name: %s: %s", category, i.Name()))
		}
		dupeChecker[i.Name()] = true
		for p := range i.Prerequisites() {
			if _, ok := t.allProperties[p]; !ok {
				panic(fmt.Sprintf("Prerequisite property not found in allProperties, is it registered?: %s", p))
			}
		}
	}
	t.AllItems[category] = items
}

func (r *StructRegistry) HasProperty(p Property) bool {
	_, ok := r.allProperties[p]
	return ok
}
