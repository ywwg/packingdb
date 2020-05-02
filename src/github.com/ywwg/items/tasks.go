package items

import (
	plib "github.com/ywwg/packinglib"
)

var tasks = []plib.Item{
	plib.NewBasicItem("Cash Money", nil, nil),
	plib.NewBasicItem("permetherin", []string{"Camping", "Hiking"}, nil),
	plib.NewBasicItem("label things", []string{"Camping"}, nil),
	plib.NewBasicItem("chop broccoli thumb-size", []string{"Firefly2017"}, nil),
	plib.NewBasicItem("chop cauliflower thumb-size", []string{"Firefly2017"}, nil),
	plib.NewBasicItem("peel cut carrots diagonal disks", []string{"Firefly2017"}, nil),
	plib.NewBasicItem("snap ends off snow peas", []string{"Firefly2017"}, nil),
	plib.NewBasicItem("fix nevermore sign", []string{"Firefly2017"}, nil),
	plib.NewBasicItem("fill camelbak", []string{"Camping", "Hiking"}, nil),
	plib.NewBasicItem("check calendar TODO", []string{"Tiny House"}, nil),
	plib.NewBasicItem("CC travel alerts", []string{"BigTrip", "International"}, nil),
	plib.NewBasicItem("intl data", []string{"International"}, nil),
	plib.NewBasicItem("sync music", nil, nil),
	plib.NewBasicItem("thermostat on hold", nil, nil),
}

func init() {
	plib.RegisterItems("Tasks", tasks)
}
