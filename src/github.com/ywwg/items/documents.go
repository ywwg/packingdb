package items

import (
	plib "github.com/ywwg/packinglib"
)

var documents = []plib.Item{
	plib.NewBasicItem("Wallet", nil, nil),
	plib.NewBasicItem("Passport", []string{"International"}, nil),
	plib.NewBasicItem("Ticket", []string{"PaidEvent"}, nil),
}

func init() {
	plib.RegisterItems("Documents", documents)
}
