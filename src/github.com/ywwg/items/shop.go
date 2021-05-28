package items

import (
	plib "github.com/ywwg/packinglib"
)

var toBuy = []plib.Item{}

func init() {
	plib.RegisterItems("To Buy", toBuy)
}
