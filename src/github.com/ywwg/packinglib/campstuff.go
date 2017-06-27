package packinglib

var campStuff = []Item{
	NewBasicItem("tent", PropertySet{Camping: true}),
	NewBasicItem("sleeping bag", PropertySet{Camping: true}),
	NewBasicItem("sleeping pad", PropertySet{Camping: true}),
	NewBasicItem("tent light", PropertySet{Camping: true}),
	NewBasicItem("headlamp", PropertySet{Dark: true}),
	NewBasicItem("towel for butt driving home", PropertySet{Burn: true}),
	NewBasicItem("tarps", PropertySet{Burn: true}),
	NewBasicItem("ropes", PropertySet{Burn: true}),
	NewBasicItem("lighter/matches", PropertySet{Camping: true}),
	NewBasicItem("bug spray", PropertySet{Camping: true}),
	NewBasicItem("dirty clothes bag", PropertySet{Camping: true}),
	NewConsumableItem("trash bag", 0.25, NoUnits, PropertySet{Camping: true}),
	NewConsumableItem("recycle bag", 0.125, NoUnits, PropertySet{Camping: true}),
	NewConsumableItem("drinking water", 0.5, "gallons", PropertySet{Burn: true}),
	NewConsumableItem("cooking water", 0.25, "gallons", PropertySet{Burn: true}),
	NewBasicItem("camelbak", PropertySet{Burn: true, Camping: true}),
	NewBasicItem("cart", PropertySet{Burn: true}),
	NewBasicItem("TASK: permetherin", PropertySet{Camping: true}),
}

func init() {
	RegisterItems("Camping Stuff", campStuff)
}
