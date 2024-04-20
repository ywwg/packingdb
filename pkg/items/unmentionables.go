package items

import (
	plib "github.com/ywwg/packingdb/pkg/packinglib"
)

var unmentionables = []*plib.Item{
	plib.NewItem("con-dams", nil, nil).Consumable(1.0, plib.NoUnits),
	plib.NewItem("lube", nil, nil),
	plib.NewItem("trees", []string{"Burn", "Camping", "Con", "Tiny House", "NYC"}, []string{"Flight"}),
	plib.NewItem("tape head cleaner", []string{"Burn", "Camping", "Con", "Tiny House", "NYC"}, []string{"Flight"}),
	plib.NewItem("pretrip shave", nil, nil),
	plib.NewItem("sleep meds", nil, nil),
	plib.NewItem("diprolene", nil, nil),
	plib.NewItem("late nite viewing", nil, nil).Consumable(1.0, "videos").Max(4.0),
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
	plib.NewItem("head", []string{"Suiting"}, nil),
	plib.NewItem("stuffing sock", []string{"Suiting"}, nil),
	plib.NewItem("paws", []string{"Suiting"}, nil),
	plib.NewItem("arm sleeves", []string{"Suiting"}, nil),
	plib.NewItem("arm sleeve elastic thing", []string{"Suiting"}, nil),
	plib.NewItem("tail", []string{"Suiting", "Tiny House"}, nil),
	plib.NewItem("badges", []string{"fancon"}, nil),
	plib.NewItem("lanyard", []string{"Con"}, nil),
	plib.NewItem("cool balaklava", []string{"Suiting"}, nil),
	plib.NewItem("in-ear earbuds", []string{"fancon"}, nil),
	plib.NewItem("games", []string{"Con"}, nil),
	plib.NewItem("sewing kit", []string{"Suiting", "BigTrip"}, nil),
	plib.NewItem("grooming brush", []string{"Suiting"}, nil),
	plib.NewItem("lint roller", []string{"Suiting"}, nil),
	plib.NewItem("glue gun", []string{"Suiting"}, nil),
	plib.NewItem("dryer sheets", []string{"Suiting"}, nil),
	plib.NewItem("suiting outfit", []string{"Suiting"}, nil),
	plib.NewItem("fake headphones", []string{"Suiting"}, nil),
}

func init() {
	plib.RegisterItems("Unmentionables", unmentionables)
	plib.RegisterItems("Suiting", suiting)
	plib.RegisterItems("Playtime", playtime)
}
