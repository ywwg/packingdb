package items

import (
	plib "github.com/ywwg/packingdb/pkg/packinglib"
)

var art = []plib.Item{
	plib.NewBasicItem("sketchbook", []string{"Art"}, nil),
	plib.NewBasicItem("pencil box", []string{"Art"}, nil),
	plib.NewBasicItem("ipad", []string{"Art"}, nil),
	plib.NewBasicItem("apple pencil coupler", []string{"Art"}, nil),
	plib.NewBasicItem("ipad cable", []string{"Art"}, nil),
}

func init() {
	plib.RegisterItems("Art", art)
}
