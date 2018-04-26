package items

import (
	plib "github.com/ywwg/packinglib"
)

var art = []plib.Item{
	plib.NewBasicItem("sketchbook", []string{"Art"}, nil),
	plib.NewBasicItem("pencil box", []string{"Art"}, nil),
	plib.NewBasicItem("copics", []string{"Art"}, nil),
}

func init() {
	plib.RegisterItems("Art", art)
}
