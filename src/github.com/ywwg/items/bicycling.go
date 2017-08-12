package items

import (
	plib "github.com/ywwg/packinglib"
)

var bicycling = []plib.Item{
	plib.NewBasicItem("ankle straps", []string{"cycling"}, nil),
}

func init() {
	plib.RegisterItems("Bicycling", bicycling)
}
