package packinglib

import (
	"fmt"
	"math"
)

const NoUnits = "nounits"

type Item interface {
	// Name returns the name of the item
	Name() string

	// Satisfies returns true if the item belongs in the given context
	Satisfies(*Context) bool

	// Itemize tells the item to calculate how much of itself is needed given the context and returns the item
	Itemize(*Trip) Item

	// Count returns the number of this item that got packed
	Count() float64

	// String prints out a string representation of the packed item(s)
	String() string

	// Pack set the packed value to true
	Pack()

	// Packed returns true if the item has been packed
	Packed() bool
}

// BasicItem is the simplest item -- just prerequisites and no count, like "tent"
type BasicItem struct {
	// Name of the item.
	name string

	// count is the number of this thing that should get packed.
	count float64

	// packed is true if the item has been packed.
	packed bool

	// Prerequisites is a set of all properties that the context must have for this item to appear.
	Prerequisites PropertySet
}

func NewBasicItem(name string, allow, disallow []string) *BasicItem {
	return &BasicItem{
		name:          name,
		Prerequisites: buildPropertySet(allow, disallow),
	}
}

func (i *BasicItem) Name() string {
	return i.name
}

// Satisfies returns true if the context satisfies the item's requirements.
func (i *BasicItem) Satisfies(c *Context) bool {
	// Any property satisfies (OR)
	if len(i.Prerequisites) == 0 {
		return true
	}
	found := false
	// If all prereqs are denies, we can return true as long as none of the
	// denials were activated (no need for a positive requirement).
	allDenies := true
	for p, allow := range i.Prerequisites {
		if allow == true {
			allDenies = false
		}
		// Any item that has a disallowing prerequisite immediately dissatisfies.
		if _, ok := c.Properties[p]; ok {
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

func (i *BasicItem) Itemize(t *Trip) Item {
	p := &BasicItem{}
	*p = *i
	if p.Satisfies(t.C) {
		p.count = 1.0
	}
	return p
}

func (i *BasicItem) Count() float64 {
	return i.count
}

func (i *BasicItem) String() string {
	checkbox := "☐"
	if i.packed {
		checkbox = "☑"
	}
	return fmt.Sprintf("%s %s", checkbox, i.name)
}

func (i *BasicItem) Pack() {
	i.packed = true
}

func (i *BasicItem) Packed() bool {
	return i.packed
}

type TemperatureItem struct {
	BasicItem

	// TemperatureMin is the anticipated minimum temperature.
	TemperatureMin int

	// TemperatureMax is the anticipated maximum temperature.
	TemperatureMax int
}

func NewTemperatureItem(name string, min, max int, allow, disallow []string) *TemperatureItem {
	return &TemperatureItem{
		BasicItem:      *NewBasicItem(name, allow, disallow),
		TemperatureMin: min,
		TemperatureMax: max,
	}
}

func (i *TemperatureItem) Satisfies(c *Context) bool {
	if i.TemperatureMax < c.TemperatureMin {
		return false
	}
	if i.TemperatureMin > c.TemperatureMax {
		return false
	}

	return i.BasicItem.Satisfies(c)
}

func (i *TemperatureItem) Itemize(t *Trip) Item {
	p := &TemperatureItem{}
	*p = *i
	if p.Satisfies(t.C) {
		p.count = 1.0
	}
	return p
}

type ConsumableItem struct {
	BasicItem

	// DailyRate is how much the thing gets used per day.
	DailyRate float64

	// What units the rate is in.  Use NoUnits for things without "of" qualifiers. ("1 car")
	Units string

	// Prerequisites is a set of all properties that the context must have for this item to appear.
	Prerequisites map[Property]bool
}

func NewConsumableItem(name string, rate float64, units string, allow, disallow []string) *ConsumableItem {
	return &ConsumableItem{
		BasicItem: *NewBasicItem(name, allow, disallow),
		DailyRate: rate,
		Units:     units,
	}
}

func (i *ConsumableItem) Itemize(t *Trip) Item {
	p := &ConsumableItem{}
	*p = *i
	if p.Satisfies(t.C) {
		p.count = math.Ceil(i.DailyRate * float64(t.Nights))
	}
	return p
}

func (i *ConsumableItem) String() string {
	checkbox := "☐"
	if i.packed {
		checkbox = "☑"
	}
	if i.Units == NoUnits {
		if i.count == float64(int(i.count)) {
			return fmt.Sprintf("%s %d %s", checkbox, int(i.count), i.name)
		} else {
			return fmt.Sprintf("%s %.1f %s", checkbox, i.count, i.name)
		}
	} else {
		if i.count == float64(int(i.count)) {
			return fmt.Sprintf("%s %d %s of %s", checkbox, int(i.count), i.Units, i.name)
		} else {
			return fmt.Sprintf("%s %.1f %s of %s", checkbox, i.count, i.Units, i.name)
		}
	}
}

type CustomConsumableItem struct {
	ConsumableItem

	// DailyRate is how much the thing gets used per day.
	RateFunc func(nights int) float64
}

func NewCustomConsumableItem(name string, rateFunc func(nights int) float64, units string, allow, disallow []string) *CustomConsumableItem {
	return &CustomConsumableItem{
		ConsumableItem: *NewConsumableItem(name, 0, units, allow, disallow),
		RateFunc:       rateFunc,
	}
}

func (i *CustomConsumableItem) Itemize(t *Trip) Item {
	p := &CustomConsumableItem{}
	*p = *i
	if p.Satisfies(t.C) {
		p.count = i.RateFunc(t.Nights)
	}
	return p
}

type ConsumableTemperatureItem struct {
	ConsumableItem
	TemperatureItem
}

func NewConsumableTemperatureItem(name string, rate float64, units string, min, max int, allow, disallow []string) *ConsumableTemperatureItem {
	return &ConsumableTemperatureItem{
		ConsumableItem:  *NewConsumableItem(name, rate, units, allow, disallow),
		TemperatureItem: *NewTemperatureItem(name, min, max, allow, disallow),
	}
}

func (i *ConsumableTemperatureItem) Satisfies(c *Context) bool {
	if !i.TemperatureItem.Satisfies(c) {
		return false
	}

	return i.ConsumableItem.Satisfies(c)
}

func (i *ConsumableTemperatureItem) Itemize(t *Trip) Item {
	p := &ConsumableTemperatureItem{}
	*p = *i
	if p.Satisfies(t.C) {
		p.ConsumableItem.count = math.Ceil(i.DailyRate * float64(t.Nights))
	}
	return p
}

func (i *ConsumableTemperatureItem) Name() string {
	return i.ConsumableItem.Name()
}

func (i *ConsumableTemperatureItem) Count() float64 {
	return i.ConsumableItem.count
}

func (i *ConsumableTemperatureItem) String() string {
	return i.ConsumableItem.String()
}

func (i *ConsumableTemperatureItem) Packed() bool {
	return i.ConsumableItem.Packed()
}

func (i *ConsumableTemperatureItem) Pack() {
	i.ConsumableItem.Pack()
}
