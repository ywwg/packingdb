package packinglib

import (
)

const NoUnits = "nounits"

type Item struct {
  // Name of the item.
  Name string

  // DailyRate is how much the thing gets used per day.  0 for non-consumables.
  DailyRate float64

  // What units the rate is in.  Use NoUnits for things without "of" qualifiers. ("1 car")
  Units string

  // TemperatureMin is the minimum temperature it makes sense to use/have the thing.  It's assumed
  // you're outside.
  TemperatureMin int

  // TemperatureMax is the maximum temperature it makes sense to use/have the thing.
  TemperatureMax int

  // Prerequisites is a set of all properties that the context must have for this item to appear.
  Prerequisites map[Property]bool
}

type PackedItem struct {
  // Name of the Item
  Name string

  // Count of how much you need
  Count float64

  // Units for the count
  Units string
}
