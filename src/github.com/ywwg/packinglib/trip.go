package packinglib

import (
  "math"
)

// Trip describes a trip, which includes a length and a context
type Trip struct {
  Days int
  C Context
}

// MakeList returns a slice of PackedItems for the given trip
func (t *Trip) MakeList() []PackedItem {
  AllItems := Clothing
  AllItems = append(AllItems, CampStuff...)
  AllItems = append(AllItems, WaterStuff...)
  AllItems = append(AllItems, Nevermore...)

  var packed []PackedItem
  for _, i := range AllItems {
    p := PackedItem {
      Name: i.Name,
      Count: 1,
      Units: i.Units,
    }
    if t.C.Satisfies(&i) {
      if i.DailyRate != 0 {
        p.Count = math.Ceil(i.DailyRate * float64(t.Days))
      }
      packed = append(packed, p)
    }
  }

  return packed
}
