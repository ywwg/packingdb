package packinglib

import (
	"fmt"
	"math"
)

// NoUnits is used when an item isn't counted by a unit word.
const NoUnits = "nounits"

// Item represents a thing that gets packed for a trip.
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

	// Pack set the packed value to whatever is passed in
	Pack(bool)

	// Packed returns true if the item has been packed
	Packed() bool

	// Prerequisites returns the PropertySet of prereqs for this item
	Prerequisites() PropertySet
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
	prerequisites PropertySet
}

// NewBasicItem creates a Basic Item with the provided allow and disallow property prerequisites.
func NewBasicItem(name string, allow, disallow []string) *BasicItem {
	return &BasicItem{
		name:          name,
		prerequisites: buildPropertySet(allow, disallow),
	}
}

// Name returns the name of the item
func (i *BasicItem) Name() string {
	return i.name
}

// Satisfies returns true if the context satisfies the item's requirements.
func (i *BasicItem) Satisfies(c *Context) bool {
	// Any property satisfies (OR)
	if len(i.prerequisites) == 0 {
		return true
	}
	found := false
	// If all prereqs are denies, we can return true as long as none of the
	// denials were activated (no need for a positive requirement).
	allDenies := true
	for p, allow := range i.prerequisites {
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

// Itemize tells the item to calculate how much of itself is needed given the context and returns the item
func (i *BasicItem) Itemize(t *Trip) Item {
	p := &BasicItem{}
	*p = *i
	if p.Satisfies(t.C) {
		p.count = 1.0
	}
	return p
}

// Count returns the number of this item should be packed.
func (i *BasicItem) Count() float64 {
	return i.count
}

// String constructs a pretty string for printing this item, including a checkbox
// for its packed status
func (i *BasicItem) String() string {
	checkbox := "☐"
	if i.packed {
		checkbox = "☑"
	}
	return fmt.Sprintf("%s %s", checkbox, i.name)
}

// Pack logs the item as packed.
func (i *BasicItem) Pack(p bool) {
	i.packed = p
}

// Packed returns true if the item has been packed
func (i *BasicItem) Packed() bool {
	return i.packed
}

// Prerequisites returns the PropertySet of prereqs for this item
func (i *BasicItem) Prerequisites() PropertySet {
	return i.prerequisites
}

// TemperatureItem represents an item that only applies in a certain temperature range.
type TemperatureItem struct {
	BasicItem

	// TemperatureMin is the anticipated minimum temperature.
	TemperatureMin int

	// TemperatureMax is the anticipated maximum temperature.
	TemperatureMax int
}

// NewTemperatureItem constructs an item given the temperature range and prereqs.
func NewTemperatureItem(name string, min, max int, allow, disallow []string) *TemperatureItem {
	return &TemperatureItem{
		BasicItem:      *NewBasicItem(name, allow, disallow),
		TemperatureMin: min,
		TemperatureMax: max,
	}
}

// Satisfies returns true if the context satisfies the item's requirements.
func (i *TemperatureItem) Satisfies(c *Context) bool {
	if i.TemperatureMax < c.TemperatureMin {
		return false
	}
	if i.TemperatureMin > c.TemperatureMax {
		return false
	}

	return i.BasicItem.Satisfies(c)
}

// Itemize tells the item to calculate how much of itself is needed given the context and returns the item
func (i *TemperatureItem) Itemize(t *Trip) Item {
	p := &TemperatureItem{}
	*p = *i
	if p.Satisfies(t.C) {
		p.count = 1.0
	}
	return p
}

// ConsumableItem is an item where a certain number will be used every day.
type ConsumableItem struct {
	BasicItem

	// DailyRate is how much the thing gets used per day.
	DailyRate float64

	// What units the rate is in.  Use NoUnits for things without "of" qualifiers. ("1 car")
	Units string

	// prerequisites is a set of all properties that the context must have for this item to appear.
	prerequisites map[Property]bool
}

// NewConsumableItem constructs an item with the given rate of usage and other prereqs.
func NewConsumableItem(name string, rate float64, units string, allow, disallow []string) *ConsumableItem {
	return &ConsumableItem{
		BasicItem: *NewBasicItem(name, allow, disallow),
		DailyRate: rate,
		Units:     units,
	}
}

// Itemize tells the item to calculate how much of itself is needed given the context and returns the item
func (i *ConsumableItem) Itemize(t *Trip) Item {
	p := &ConsumableItem{}
	*p = *i
	if p.Satisfies(t.C) {
		p.count = math.Ceil(i.DailyRate * float64(t.Nights))
	}
	return p
}

// String constructs a pretty string for printing this item, including a checkbox
// for its packed status
func (i *ConsumableItem) String() string {
	checkbox := "☐"
	if i.packed {
		checkbox = "☑"
	}
	if i.Units == NoUnits {
		if i.count == float64(int(i.count)) {
			return fmt.Sprintf("%s %d %s", checkbox, int(i.count), i.name)
		}
		return fmt.Sprintf("%s %.1f %s", checkbox, i.count, i.name)
	}
	if i.count == float64(int(i.count)) {
		return fmt.Sprintf("%s %d %s of %s", checkbox, int(i.count), i.Units, i.name)
	}
	return fmt.Sprintf("%s %.1f %s of %s", checkbox, i.count, i.Units, i.name)
}

// ConsumableMaxItem represents an item that you may need multiple of, but at some
// point it maxes out and there's no point bringing more.
type ConsumableMaxItem struct {
	ConsumableItem

	// Max is the most of these you'll ever need.
	Max float64
}

// NewConsumableMaxItem constructs an item with the given rate of usage, maximum, and other prereqs.
func NewConsumableMaxItem(name string, rate float64, max float64, units string, allow, disallow []string) *ConsumableMaxItem {
	return &ConsumableMaxItem{
		ConsumableItem: *NewConsumableItem(name, rate, units, allow, disallow),
		Max:            max,
	}
}

// Itemize tells the item to calculate how much of itself is needed given the context and returns the item
func (i *ConsumableMaxItem) Itemize(t *Trip) Item {
	p := &ConsumableMaxItem{}
	*p = *i
	if p.Satisfies(t.C) {
		p.count = math.Min(math.Ceil(i.DailyRate*float64(t.Nights)), p.Max)
	}
	return p
}

// CustomConsumableItem is a consumable item that takes a function to determine how many
// are needed, instead of a simple float rate.
type CustomConsumableItem struct {
	ConsumableItem

	// DailyRate is how much the thing gets used per day.
	RateFunc func(nights int, props PropertySet) float64
}

// NewCustomConsumableItem constructs an item with the given rate function and other prereqs.
func NewCustomConsumableItem(name string, rateFunc func(nights int, props PropertySet) float64, units string, allow, disallow []string) *CustomConsumableItem {
	return &CustomConsumableItem{
		ConsumableItem: *NewConsumableItem(name, 0, units, allow, disallow),
		RateFunc:       rateFunc,
	}
}

// Itemize tells the item to calculate how much of itself is needed given the context and returns the item
func (i *CustomConsumableItem) Itemize(t *Trip) Item {
	p := &CustomConsumableItem{}
	*p = *i
	if p.Satisfies(t.C) {
		p.count = i.RateFunc(t.Nights, t.C.Properties)
	}
	return p
}

// ConsumableTemperatureItem is both consumable and only applies for a given temperature range.
// Diamond pattern in action!
type ConsumableTemperatureItem struct {
	ConsumableItem
	TemperatureItem
}

// NewConsumableTemperatureItem constructs the item with the given settings.
func NewConsumableTemperatureItem(name string, rate float64, units string, min, max int, allow, disallow []string) *ConsumableTemperatureItem {
	return &ConsumableTemperatureItem{
		ConsumableItem:  *NewConsumableItem(name, rate, units, allow, disallow),
		TemperatureItem: *NewTemperatureItem(name, min, max, allow, disallow),
	}
}

// Satisfies returns true if the item belongs in the given context
func (i *ConsumableTemperatureItem) Satisfies(c *Context) bool {
	if !i.TemperatureItem.Satisfies(c) {
		return false
	}

	return i.ConsumableItem.Satisfies(c)
}

// Itemize tells the item to calculate how much of itself is needed given the context and returns the item
func (i *ConsumableTemperatureItem) Itemize(t *Trip) Item {
	p := &ConsumableTemperatureItem{}
	*p = *i
	if p.Satisfies(t.C) {
		p.ConsumableItem.count = math.Ceil(i.DailyRate * float64(t.Nights))
	}
	return p
}

// Name returns the name of the item
func (i *ConsumableTemperatureItem) Name() string {
	return i.ConsumableItem.Name()
}

// Count returns the number of this item that got packed
func (i *ConsumableTemperatureItem) Count() float64 {
	return i.ConsumableItem.count
}

// String constructs a pretty string for printing this item, including a checkbox
// for its packed status
func (i *ConsumableTemperatureItem) String() string {
	return i.ConsumableItem.String()
}

// Packed returns true if the item has been packed
func (i *ConsumableTemperatureItem) Packed() bool {
	return i.ConsumableItem.Packed()
}

// Pack logs the item as packed.
func (i *ConsumableTemperatureItem) Pack(p bool) {
	i.ConsumableItem.Pack(p)
}

// Prerequisites returns the PropertySet of prereqs for this item
func (i *ConsumableTemperatureItem) Prerequisites() PropertySet {
	return i.ConsumableItem.Prerequisites()
}

// ConsumableMaxTemperatureItem is both consumable, with a maximum amount, and has a temperature range.
// Diamond pattern again!
type ConsumableMaxTemperatureItem struct {
	ConsumableMaxItem
	TemperatureItem
}

// NewConsumableMaxTemperatureItem constructs the item with the many varied requirements.
func NewConsumableMaxTemperatureItem(name string, rate float64, maxNum float64, units string, min, max int, allow, disallow []string) *ConsumableMaxTemperatureItem {
	return &ConsumableMaxTemperatureItem{
		ConsumableMaxItem: *NewConsumableMaxItem(name, rate, maxNum, units, allow, disallow),
		TemperatureItem:   *NewTemperatureItem(name, min, max, allow, disallow),
	}
}

// Satisfies returns true if the item belongs in the given context
func (i *ConsumableMaxTemperatureItem) Satisfies(c *Context) bool {
	if !i.TemperatureItem.Satisfies(c) {
		return false
	}

	return i.ConsumableMaxItem.Satisfies(c)
}

// Itemize tells the item to calculate how much of itself is needed given the context and returns the item
func (i *ConsumableMaxTemperatureItem) Itemize(t *Trip) Item {
	p := &ConsumableMaxTemperatureItem{}
	*p = *i
	if p.Satisfies(t.C) {
		p.ConsumableMaxItem.count = math.Min(math.Ceil(i.DailyRate*float64(t.Nights)), p.Max)
	}
	return p
}

// Name returns the name of the item
func (i *ConsumableMaxTemperatureItem) Name() string {
	return i.ConsumableMaxItem.Name()
}

// Count returns the number of this item that got packed
func (i *ConsumableMaxTemperatureItem) Count() float64 {
	return i.ConsumableMaxItem.count
}

// String constructs a pretty string for printing this item, including a checkbox
// for its packed status
func (i *ConsumableMaxTemperatureItem) String() string {
	return i.ConsumableMaxItem.String()
}

// Packed returns true if the item has been packed
func (i *ConsumableMaxTemperatureItem) Packed() bool {
	return i.ConsumableMaxItem.Packed()
}

// Pack logs the item as packed.
func (i *ConsumableMaxTemperatureItem) Pack(p bool) {
	i.ConsumableMaxItem.Pack(p)
}
