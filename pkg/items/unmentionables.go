package items

import (
	plib "github.com/ywwg/packingdb/pkg/packinglib"
)

var unmentionables = []*plib.Item{
	plib.NewItem("con-dams", nil, nil).Consumable(1.0),
	plib.NewItem("lube", nil, nil),
	plib.NewItem("trees", []string{"Burn", "Camping", "Con", "Tiny House", "NYC"}, []string{"Flight"}),
	plib.NewItem("tape head cleaner", []string{"Burn", "Camping", "Con", "Tiny House", "NYC"}, []string{"Flight"}),
	plib.NewItem("pretrip shave", nil, nil),
	plib.NewItem("sleep meds", nil, nil),
	plib.NewItem("diprolene", nil, nil),
	plib.NewItem("late nite viewing", nil, nil).Units("videos").Consumable(1.0).Max(4.0),
	plib.NewItem("late nite viewing HD", []string{"Tiny House"}, nil),
}

var playtime = []*plib.Item{
	plib.NewItem("harness", []string{"Playtime"}, nil),
	plib.NewItem("collar and leash", []string{"Playtime"}, nil),
	plib.NewItem("poppers", []string{"Playtime"}, []string{"Flight"}),
	plib.NewItem("ring vibe", []string{"Playtime"}, nil),
	plib.NewItem("dildoes", []string{"Playtime"}, nil),
	plib.NewItem("toy cleaner", []string{"Playtime"}, nil),
}

var suiting = []*plib.Item{
	plib.NewItem("heads", []string{"Suiting"}, nil),
	plib.NewItem("bodysuit", []string{"Suiting"}, nil),
	plib.NewItem("feet", []string{"Suiting"}, nil),
	plib.NewItem("stuffing sock", []string{"Suiting"}, nil),
	plib.NewItem("paws", []string{"Suiting"}, nil),
	plib.NewItem("arm sleeves", []string{"Suiting"}, nil),
	plib.NewItem("arm sleeve elastic thing", []string{"Suiting"}, nil),
	plib.NewItem("tail", []string{"Suiting", "Tiny House"}, nil),
	plib.NewItem("badges", []string{"fancon"}, nil),
	plib.NewItem("charged LED badge", []string{"fancon"}, nil),
	plib.NewItem("lanyard", []string{"Con"}, nil),
	plib.NewItem("Midnight Makers lanyard", []string{"Con"}, nil),
	plib.NewItem("cool balaklava", []string{"Suiting"}, nil),
	plib.NewItem("cooling top", []string{"Suiting"}, nil),
	plib.NewItem("cooling bottom", []string{"Suiting"}, nil),
	plib.NewItem("head battery pack", []string{"Suiting"}, nil),
	plib.NewItem("in-ear earbuds", []string{"fancon"}, nil),
	plib.NewItem("daki", []string{"fancon"}, nil),
	plib.NewItem("extra blanket", []string{"fancon"}, nil),
	plib.NewItem("games", []string{"Con"}, nil),
	plib.NewItem("sewing kit", []string{"Suiting", "BigTrip"}, nil),
	plib.NewItem("grooming brush", []string{"Suiting"}, nil),
	plib.NewItem("lint roller", []string{"Suiting"}, nil),
	plib.NewItem("glue gun", []string{"Suiting"}, nil),
	plib.NewItem("dryer sheets", []string{"Suiting"}, nil),
	plib.NewItem("suiting outfit", []string{"Suiting"}, nil),
	plib.NewItem("Huxley/Laelia stickers", []string{"fancon"}, nil),
}
