package packinglib

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// ErrContextExists is returned by RegisterContext when a context with the
// given name is already registered. Callers that expect to re-register (e.g.,
// reload paths) should use GetContext or GetConcreteContext instead.
var ErrContextExists = errors.New("context already registered")

// ContextRegistry is a generic interface for uhhh doing things with contexts.
type Registry interface {
	ContextList() []string
	AllItems() PackList
	AllProperties() map[Property]string
	GetDescription(p Property) string
	Context(name string) (*Context, error)

	// Not clear how registration should work... probably just "add"
	RegisterContext(c Context) error
	RegisterProperty(prop Property, desc string)
	RegisterItems(category Category, items []*Item)
	ListProperties() []Property
	GetContext(name string) (*Context, error)
	GetConcreteContext(name string, nights, tmin, tmax int) (*Context, error)
	HasProperty(p Property) bool
}

// StructRegistry is a simple in-memory impl of a ContextRegistry
type StructRegistry struct {
	contexts      map[string]Context
	allProperties map[Property]string
	allItems      PackList
	dupeChecker   map[string]bool
}

func NewStructRegistry() Registry {
	r := &StructRegistry{
		contexts:      make(map[string]Context),
		allProperties: make(map[Property]string),
		allItems:      make(PackList),
		dupeChecker:   make(map[string]bool),
	}
	return r
}

// ContextList returns a sorted slice of strings of the contexts.
func (r *StructRegistry) ContextList() []string {
	keys := make([]string, len(r.contexts))
	i := 0
	for k := range r.contexts {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func (r *StructRegistry) Context(name string) (*Context, error) {
	c, ok := r.contexts[name]
	if !ok {
		return nil, fmt.Errorf("unknown context: %s", name)
	}
	return &c, nil
}

func (r *StructRegistry) AllItems() PackList {
	return r.allItems
}

func (r *StructRegistry) AllProperties() map[Property]string {
	return r.allProperties
}

func (r *StructRegistry) GetDescription(p Property) string {
	if desc, ok := r.allProperties[p]; ok {
		return desc
	}
	return "no description what now"
}

// RegisterContext registers the given context with the system.
// Also registers a property with the context name.
// Returns ErrContextExists wrapped with the name if the context is already registered.
func (r *StructRegistry) RegisterContext(c Context) error {
	if _, ok := r.contexts[c.Name]; ok {
		return fmt.Errorf("%w: %s", ErrContextExists, c.Name)
	}
	r.contexts[c.Name] = c
	r.RegisterProperty(Property(c.Name), "")
	return nil
}

// RegisterProperty adds a new Property to the database so it can be used.
// desc should be a user-visible description of the property.
// Does not verify that all of the properties are in the allProperties map.
func (r *StructRegistry) RegisterProperty(prop Property, desc string) {
	if _, ok := r.allProperties[prop]; ok {
		if desc != "" {
			r.allProperties[prop] = desc
		}
		return
	}
	r.allProperties[prop] = desc
}

// GetContext returns the context of the given name, or returns error if not found.
func (r *StructRegistry) GetContext(name string) (*Context, error) {
	c := &Context{}
	found, ok := r.contexts[name]
	if !ok {
		return nil, fmt.Errorf("unknown context: %s", name)
	}
	*c = found
	return c, nil
}

// GetConcreteContext loads the given context and substitutes the provided
// temperature range.
func (r *StructRegistry) GetConcreteContext(name string, nights, tmin, tmax int) (*Context, error) {
	c, err := r.GetContext(name)
	if err != nil {
		return nil, err
	}
	c.Nights = nights
	c.TemperatureMin = tmin
	c.TemperatureMax = tmax
	return c, nil
}

// RegisterItems appends the given slice of Items to the registry under
// the given category.  Duplicate categories will be appended.  Items
// with duplicate case-insensitive names, even across categories,
// cause a panic.
func (r *StructRegistry) RegisterItems(category Category, items []*Item) {
	for _, i := range items {
		if _, ok := r.dupeChecker[strings.ToLower(i.Name())]; ok {
			panic(fmt.Sprintf("Duplicate item name: %s: %s", category, i.Name()))
		}
		r.dupeChecker[i.Name()] = true
		for p := range i.Prerequisites() {
			if _, ok := r.allProperties[p]; !ok {
				panic(fmt.Sprintf("Prerequisite property not found in allProperties, is it registered?: %s", p))
			}
		}
		// Pre-sort mutators by Priority() descending, once at registration time.
		// After this point mutator order is invariant. AdjustCount relies on this
		// and must not re-sort (see ISOL-07 / D-05).
		sort.Slice(i.mutators, func(x, y int) bool {
			return i.mutators[y].Priority() < i.mutators[x].Priority()
		})
	}
	r.allItems[category] = items
}

// ListProperties returns all of the registered properties as a slice of strings.
func (r *StructRegistry) ListProperties() []Property {
	var l []Property
	for k := range r.AllProperties() {
		l = append(l, k)
	}
	less := func(i, j int) bool {
		return strings.ToLower(string(l[i])) < strings.ToLower(string(l[j]))
	}
	sort.Slice(l, less)
	return l
}

func (r *StructRegistry) HasProperty(p Property) bool {
	_, ok := r.allProperties[p]
	return ok
}
