package packinglib

var Electrical = []Item{
	NewBasicItem("Generator", PropertySet{Burn: true}),
	NewConsumableItem("Gas", 1.0, "gallon", PropertySet{Burn: true}),
	NewBasicItem("Extension Cords", PropertySet{Burn: true}),
	NewBasicItem("Cell Phone Charger", nil),
}
