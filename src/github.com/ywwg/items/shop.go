package items

import (
	plib "github.com/ywwg/packinglib"
)

var toBuy = []plib.Item{}

func init() {
	plib.RegisterProperty("Firefly2017")
	plib.RegisterItems("To Buy", toBuy)
}
