package items

import (
	plib "github.com/ywwg/packinglib"
)

var pride = []plib.Item{
	plib.NewBasicItem("pride shirt", []string{"PortlandPride"}, nil),
}

var tubing2021 = []plib.Item{
	plib.NewBasicItem("big pot", []string{"Tubing"}, nil),
	plib.NewBasicItem("bottle opener", []string{"Tubing"}, nil),
	plib.NewBasicItem("chef knife", []string{"Tubing"}, nil),
	plib.NewBasicItem("float cooler", []string{"Tubing"}, nil),
	plib.NewBasicItem("salt / pepper", []string{"Tubing"}, nil),
	plib.NewBasicItem("serving spoon", []string{"Tubing"}, nil),
	plib.NewBasicItem("syringes", []string{"Tubing"}, nil),
}

func init() {
	plib.RegisterProperty(plib.Property("PortlandPride"), "Portland Pride Parade!")
	plib.RegisterItems("PortlandPride", pride)
}
