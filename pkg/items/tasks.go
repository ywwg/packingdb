package items

import (
	plib "github.com/ywwg/packingdb/pkg/packinglib"
)

var tasks = []*plib.Item{
	plib.NewItem("Cash Money", nil, nil),
	plib.NewItem("permetherin", []string{"Camping", "Hiking"}, nil),
	plib.NewItem("label things", []string{"Camping"}, nil),
	plib.NewItem("fill camelbak", []string{"Camping", "Hiking"}, nil),
	plib.NewItem("check calendar TODO", []string{"Tiny House"}, nil),
	plib.NewItem("CC travel alerts", []string{"BigTrip", "International"}, nil),
	plib.NewItem("intl data", []string{"International"}, nil),
	plib.NewItem("sync music", nil, nil),
	plib.NewItem("thermostat on hold", nil, nil),
}

func init() {
	plib.RegisterItems("Tasks", tasks)
}
