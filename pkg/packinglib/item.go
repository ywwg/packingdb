package packinglib

import (
	"math"
)

// NoUnits is used when an item isn't counted by a unit word.
const NoUnits = "nounits"

// Item
type Item struct {
	// Name of the item.
	name string

	// count is the number of this thing that should get packed.
	count float64

	// packed is true if the item has been packed.
	packed bool

	// Prerequisites is a set of all properties that the context must have for this item to appear.
	prerequisites PropertySet

	mutators []packMutator
}

// NewItem creates a Basic Item with the provided allow and disallow property prerequisites.
func NewItem(name string, allow, disallow []string) *Item {
	return &Item{
		name:          name,
		prerequisites: buildPropertySet(allow, disallow),
	}
}

// Name returns the name of the item
func (i *Item) Name() string {
	return i.name
}

// Satisfies returns true if the context satisfies the item's requirements.
func (i *Item) Satisfies(c *Context) bool {
	// Any property satisfies (OR)
	if len(i.prerequisites) == 0 {
		return true
	}
	found := false
	// If all prereqs are denies, we can return true as long as none of the
	// denials were activated (no need for a positive requirement).
	allDenies := true
	for p, allow := range i.prerequisites {
		if allow {
			allDenies = false
		}
		// Any item that has a disallowing prerequisite immediately dissatisfies.
		if c.Properties[p] {
			if !allow {
				return false
			}
			found = true
		}
	}
	if allDenies {
		return true
	}
	return found
}

// AdjustCount tells the item to calculate how much of itself is needed given the
// context and returns the item. Mutators multiply together I guess???
func (i *Item) AdjustCount(c *Context) {
	i.count = 1.0
	// First check is always Satisfies
	if !i.Satisfies(c) {
		i.count = 0.0
		return
	}
	for _, t := range i.mutators {
		// this makes it feel like satisfies should be a mutator. how do we deal with
		// iterating through... maybe any time something returns zero we stop
		// processing?  i.e., any mutator can only adjust a number from 1.0, never
		// down to zero unless it means
		i.count = t.AdjustCount(c, i.count)
		if i.count == 0.0 {
			return
		}
	}
}

// Count returns the number of this item should be packed.
func (i *Item) Count() float64 {
	return i.count
}

// String constructs a pretty string for printing this item, including a checkbox
// for its packed status
// XXXX yikes we should not be decorating here
func (i *Item) String() string {
	// checkbox := "○"
	// if i.packed {
	// 	checkbox = "●"
	// }
	// return fmt.Sprintf("%s %s", checkbox, i.name)
	return i.name
}

// Pack logs the item as packed.
func (i *Item) Pack(p bool) {
	i.packed = p
}

// Packed returns true if the item has been packed
func (i *Item) Packed() bool {
	return i.packed
}

// Prerequisites returns the PropertySet of prereqs for this item
func (i *Item) Prerequisites() PropertySet {
	return i.prerequisites
}

type packMutator interface {
	// AdjustCount takes a count and adjusts it for what this mutator does. If the
	// mutator has certain requirements, it should adjust the count to 0.
	AdjustCount(c *Context, count float64) float64
}

// temperatureMutator represents an item that only applies in a certain temperature range.
type temperatureMutator struct {
	// TemperatureMin is the anticipated minimum temperature.
	TemperatureMin int

	// TemperatureMax is the anticipated maximum temperature.
	TemperatureMax int
}

func (i *Item) TemperatureRange(tMin, tMax int) *Item {
	i.mutators = append(i.mutators, &temperatureMutator{tMin, tMax})
	return i
}

func (i *temperatureMutator) AdjustCount(c *Context, count float64) float64 {
	if i.TemperatureMax < c.TemperatureMin {
		return 0.0
	}
	if i.TemperatureMin > c.TemperatureMax {
		return 0.0
	}

	return 1.0
}

// ConsumableItem is an item where a certain number will be used every day.
type consumableMutator struct {
	// DailyRate is how much the thing gets used per day.
	DailyRate float64

	// What units the rate is in.  Use NoUnits for things without "of" qualifiers.
	// ("1 car")
	Units string
}

func (i *Item) Consumable(rate float64, units string) *Item {
	i.mutators = append(i.mutators, &consumableMutator{rate, units})
	return i
}

// AdjustCount tells the item to calculate how much of itself is needed given
// the context and returns the item
func (i *consumableMutator) AdjustCount(c *Context, count float64) float64 {
	return count * math.Ceil(i.DailyRate*float64(c.Nights))
}

// // String constructs a pretty string for printing this item, including a checkbox
// // for its packed status
// func (i *ConsumableItem) String() string {
// 	checkbox := "○"
// 	if i.packed {
// 		checkbox = "●"
// 	}
// 	if i.Units == NoUnits {
// 		if i.count == float64(int(i.count)) {
// 			return fmt.Sprintf("%s %d %s", checkbox, int(i.count), i.name)
// 		}
// 		return fmt.Sprintf("%s %.1f %s", checkbox, i.count, i.name)
// 	}
// 	if i.count == float64(int(i.count)) {
// 		return fmt.Sprintf("%s %d %s of %s", checkbox, int(i.count), i.Units, i.name)
// 	}
// 	return fmt.Sprintf("%s %.1f %s of %s", checkbox, i.count, i.Units, i.name)
// }

// maxCountMutator represents an item that you may need multiple of, but at some
// point it maxes out and there's no point bringing more.
type maxCountMutator struct {
	// Max is the most of these you'll ever need.
	Max float64
}

func (i *Item) Max(max float64) *Item {
	i.mutators = append(i.mutators, &maxCountMutator{max})
	return i
}

// AdjustCount tells the item to calculate how much of itself is needed given the context and returns the item
func (i *maxCountMutator) AdjustCount(c *Context, count float64) float64 {
	return math.Min(count, i.Max)
}

// customConsumableMutator is a consumable item that takes a function to
// determine how many are needed, instead of a simple float rate.
type customConsumableMutator struct {
	// DailyRate is how much the thing gets used per day.
	RateFunc func(count float64, nights int, props PropertySet) float64

	Units string
}

func (i *Item) Custom(f func(count float64, nights int, props PropertySet) float64, units string) *Item {
	i.mutators = append(i.mutators, &customConsumableMutator{f, units})
	return i
}

// AdjustCount tells the item to calculate how much of itself is needed given the context and returns the item
func (i *customConsumableMutator) AdjustCount(c *Context, count float64) float64 {
	return i.RateFunc(count, c.Nights, c.Properties)
}
