package items

import (
	plib "github.com/ywwg/packinglib"
)

var tasks = []plib.Item{
	plib.NewBasicItem("permetherin", []string{"Camping"}, nil),
	plib.NewBasicItem("label things", []string{"Camping"}, nil),
}

func init() {
	plib.RegisterItems("Tasks", tasks)
}
