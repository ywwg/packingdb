package packinglib

var nevermore = []Item{
	NewBasicItem("Nevermore", PropertySet{Burn: true}),
	// TODO: like, list everything?
}

func init() {
	RegisterItems("Nevermore", nevermore)
}
