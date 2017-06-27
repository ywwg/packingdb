package packinglib

import (
)

// Property is a value describing a property the context has.  Sort of like boolean flags.
type Property int

const (
  None = iota
  Electricity
  RunningWater
  Swimming
  Dirt
  Loud
  Bright
  Dark
  Sweat
  Camping
  GrumpCamping
  Performing
  Burn
  Fancy
)

type PropertyList map[Property]bool

// Context is struct that holds data about the context of the trip
type Context struct {
  // Name of the context ("The Cape", "The Tiny House", "Firefly")
  Name string

  // TemperatureMin is the anticipated minimum temperature.
  TemperatureMin int

  // TemperatureMax is the anticipated maximum temperature.
  TemperatureMax int

  Properties PropertyList
}

// Satisfies returns true if the given item is satisfied by the context.
func (c *Context) Satisfies(i *Item) bool {
  // Temperatures don't satisfy if the temps are completely unaligned, which only happens
  // in these two cases.
  if i.TemperatureMax < c.TemperatureMin {
    return false
  }
  if i.TemperatureMin > c.TemperatureMax {
    return false
  }

  // All properties must be satisfied
  for p := range i.Prerequisites {
    if _, ok := c.Properties[p]; !ok {
      return false
    }
  }

  return true
}

var FireflyContext = Context{
  Name: "Firefly",
  TemperatureMin: 52,
  TemperatureMax: 80,
  Properties: PropertyList{
    Swimming: true,
    Dirt: true,
    Loud: true,
    Bright: true,
    Sweat: true,
    Camping: true,
    Dark: true,
    Burn: true,
    Performing: true,
  },
}
