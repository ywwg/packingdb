package items

import (
	plib "github.com/ywwg/packinglib"
)

var pride = []plib.Item{
	plib.NewBasicItem("pride shirt", []string{"PortlandPride"}, nil),
}

func init() {
	plib.RegisterProperty(plib.Property("PortlandPride"), "Portland Pride Parade!")
	plib.RegisterItems("PortlandPride", pride)
}
