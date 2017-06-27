package items

import (
	plib "github.com/ywwg/packinglib"
)

var nevermore = []plib.Item{
	plib.NewBasicItem("Nevermore", []string{"Burn"}, nil),
	// TODO: like, list everything?
}

func init() {
	plib.RegisterItems("Nevermore", nevermore)
}
