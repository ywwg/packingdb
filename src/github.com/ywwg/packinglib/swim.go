package packinglib

var WaterStuff = []Item{
	NewBasicItem("froggies", PropertySet{Swimming: true}),
	NewConsumableItem("swim suit", 0.5, NoUnits, PropertySet{Swimming: true}),
	NewBasicItem("swim towel", PropertySet{Swimming: true}),
}
