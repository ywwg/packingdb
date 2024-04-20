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

	// Register our own name as a property, this allows items to require that an
	// entire context exists (or not), which is useful.
	RegisterProperty(Property(c.Name), "")
	if err := c.addProperty(c.Name); err != nil {
		return nil, err
	}
	for _, p := range properties {
		if err := c.addProperty(p); err != nil {
			return nil, err
		}
	}

	RegisterContext(*c)
	return c, nil
}

// addProperty adds the property with the given name to the context, or returns
// error if it's not found. Empty strings are ignored.
func (c *Context) addProperty(prop string) error {
	if prop == "" {
		return nil
	}
	if _, ok := allProperties[Property(prop)]; !ok {
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
	if _, ok := allProperties[Property(prop)]; !ok {
		return fmt.Errorf("didn't find property, is it registered?: %s", prop)
	}
	c.Properties[Property(prop)] = false
	return nil
}

// hasProperty returns true if the context has this property.
func (c *Context) hasProperty(prop Property) bool {
	return c.Properties[prop]
}
