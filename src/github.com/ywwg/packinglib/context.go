package packinglib

import (
	"fmt"
)

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

// RemoveProperty removes the property with the given name to the context, or
// returns error if it's not found. Empty strings are ignored.
func (c *Context) RemoveProperty(prop string) error {
	if prop == "" {
		return nil
	}
	if _, ok := allProperties[Property(prop)]; !ok {
		return fmt.Errorf("didn't find property, is it registered?: %s", prop)
	}
	c.Properties[Property(prop)] = false
	return nil
}

// HasProperty returns true if the context has this property.
func (c *Context) HasProperty(prop string) bool {
	return c.Properties[Property(prop)]
}
