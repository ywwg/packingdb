package items

import (
	plib "github.com/ywwg/packingdb/pkg/packinglib"
)

var pride = []*plib.Item{
	plib.NewItem("pride shirt", []string{"PortlandPride"}, nil),
}

var tubing2021 = []*plib.Item{
	plib.NewItem("big pot", []string{"Tubing"}, nil),
	plib.NewItem("bottle opener", []string{"Tubing"}, nil),
	plib.NewItem("chef knife", []string{"Tubing"}, nil),
	plib.NewItem("float cooler", []string{"Tubing"}, nil),
	plib.NewItem("salt / pepper", []string{"Tubing"}, nil),
	plib.NewItem("serving spoon", []string{"Tubing"}, nil),
	plib.NewItem("syringes", []string{"Tubing"}, nil),
}

func init() {
	plib.RegisterProperty(plib.Property("PortlandPride"), "Portland Pride Parade!")
	plib.RegisterItems("PortlandPride", pride)
}
