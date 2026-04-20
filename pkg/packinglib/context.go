package packinglib

import (
	"fmt"
)

// Context is struct that holds data about the context of the trip
type Context struct {
	// Name of the context ("The Cape", "The Tiny House", "Firefly")
	Name string

	Nights int

	// TemperatureMin is the anticipated minimum temperature.
	TemperatureMin int

	// TemperatureMax is the anticipated maximum temperature.
	TemperatureMax int

	Properties PropertySet

	registry Registry
}

// NewContext creates a new context with the given name, temperature range, and properties.
// Returns nil if any of the properties is unknown.  Properties are optional.
func NewContext(r Registry, name string, nights, tmin, tmax int, properties []string) (*Context, error) {
	c := &Context{
		Name:           name,
		Nights:         nights,
		TemperatureMin: tmin,
		TemperatureMax: tmax,
		Properties:     make(PropertySet),
		registry:       r,
	}

	// Register our own name as a property, this allows items to require that an
	// entire context exists (or not), which is useful.
	r.RegisterProperty(Property(c.Name), "")
	if err := c.addProperty(c.Name); err != nil {
		return nil, err
	}
	for _, p := range properties {
		if err := c.addProperty(p); err != nil {
			return nil, err
		}
	}

	if err := r.RegisterContext(*c); err != nil {
		return nil, err
	}
	return c, nil
}

// addProperty adds the property with the given name to the context, or returns
// error if it's not found. Empty strings are ignored.
func (c *Context) addProperty(prop string) error {
	if prop == "" {
		return nil
	}
	if !c.registry.HasProperty(Property(prop)) {
		return fmt.Errorf("didn't find property, is it registered?: %s", prop)
	}
	// Recursively add contained properties that are actually also contexts. XXXXX
	// this is horrible! we should definitely not do this.
	//  Iguess we just remove this and then see what breaks... and then those things
	// should be contexts that optionally contain subcontexts for clarity. I think
	// we did this because the saved file only lists properties. Therefore we
	// should have the file save both, but disambiguate which is which.

	// That means a schema change... but that seems fine.

	// if subContext, err := GetContext(prop); err == nil {
	// 	for p, inc := range subContext.Properties {
	// 		if _, ok := c.Properties[p]; inc && !ok {
	// 			c.addProperty(string(p))
	// 		}
	// 	}
	// }
	c.Properties[Property(prop)] = true
	return nil
}

// removeProperty removes the property with the given name to the context, or
// returns error if it's not found. Empty strings are ignored.
func (c *Context) removeProperty(prop string) error {
	if prop == "" {
		return nil
	}
	if !c.registry.HasProperty(Property(prop)) {
		return fmt.Errorf("didn't find property, is it registered?: %s", prop)
	}
	c.Properties[Property(prop)] = false
	return nil
}

// hasProperty returns true if the context has this property.
func (c *Context) hasProperty(prop Property) bool {
	return c.Properties[prop]
}

// clone returns a deep copy of the context. The Properties map is copied
// into fresh backing storage so mutations on the clone do not affect the
// source. The registry reference is rebound to newRegistry so that Trip
// and Context share the same cloned registry instance (see D-03 / ISOL-04).
// Unexported: only NewTripFromCustomContext calls this. External callers
// should go through trip construction, which clones both registry and
// context consistently.
func (c *Context) clone(newRegistry Registry) *Context {
	cloned := *c // scalar fields: Name, Nights, TemperatureMin, TemperatureMax
	cloned.registry = newRegistry

	if c.Properties != nil {
		cloned.Properties = make(PropertySet, len(c.Properties))
		for k, v := range c.Properties {
			cloned.Properties[k] = v
		}
	} else {
		cloned.Properties = make(PropertySet)
	}

	return &cloned
}
