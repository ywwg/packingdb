package items

import (
	plib "github.com/ywwg/packinglib"
)

var toBuy = []plib.Item{
	plib.NewBasicItem("more rope for my own tarps?", []string{"Firefly2017"}, nil),
}

func init() {
	plib.RegisterProperty("Firefly2017")
	plib.RegisterItems("To Buy", toBuy)
}
