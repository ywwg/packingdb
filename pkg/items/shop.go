package items

import (
	plib "github.com/ywwg/packingdb/pkg/packinglib"
)

var toBuy = []plib.Item{}

func init() {
	plib.RegisterItems("To Buy", toBuy)
}
