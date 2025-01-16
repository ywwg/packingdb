package items

import (
	plib "github.com/ywwg/packinglib"
)

var unmentionables = []plib.Item{
	plib.NewConsumableItem("con-dams", 1.0, plib.NoUnits, nil, nil),
	plib.NewBasicItem("lube", nil, nil),
	plib.NewBasicItem("trees", []string{"Burn", "Camping", "Con", "Tiny House", "NYC"}, []string{"Flight"}),
	plib.NewBasicItem("tape head cleaner", []string{"Burn", "Camping", "Con", "Tiny House", "NYC"}, []string{"Flight"}),
	plib.NewBasicItem("pretrip shave", nil, nil),
	plib.NewBasicItem("sleep meds", nil, nil),
	plib.NewBasicItem("diprolene", nil, nil),
	plib.NewConsumableMaxItem("late nite viewing", 1.0, 4.0, "videos", nil, nil),
	plib.NewBasicItem("late nite viewing HD", []string{"Tiny House"}, nil),
}

var playtime = []plib.Item{
	plib.NewBasicItem("harness", []string{"Playtime"}, nil),
	plib.NewBasicItem("collar and leash", []string{"Playtime"}, nil),
	plib.NewBasicItem("poppers", []string{"Playtime"}, []string{"Flight"}),
	plib.NewBasicItem("ring vibe", []string{"Playtime"}, nil),
	plib.NewBasicItem("dildoes", []string{"Playtime"}, nil),
	plib.NewBasicItem("toy cleaner", []string{"Playtime"}, nil),
}

var suiting = []plib.Item{
	plib.NewBasicItem("head", []string{"Suiting"}, nil),
	plib.NewBasicItem("stuffing sock", []string{"Suiting"}, nil),
	plib.NewBasicItem("paws", []string{"Suiting"}, nil),
	plib.NewBasicItem("arm sleeves", []string{"Suiting"}, nil),
	plib.NewBasicItem("arm sleeve elastic thing", []string{"Suiting"}, nil),
	plib.NewBasicItem("tail", []string{"Suiting", "Tiny House"}, nil),
	plib.NewBasicItem("badges", []string{"fancon"}, nil),
	plib.NewBasicItem("lanyard", []string{"Con"}, nil),
	plib.NewBasicItem("cool balaklava", []string{"Suiting"}, nil),
	plib.NewBasicItem("cooling top", []string{"Suiting"}, nil),
	plib.NewBasicItem("cooling bottom", []string{"Suiting"}, nil),
	plib.NewBasicItem("head battery pack", []string{"Suiting"}, nil),
	plib.NewBasicItem("in-ear earbuds", []string{"fancon"}, nil),
	plib.NewBasicItem("daki", []string{"fancon"}, nil),
	plib.NewBasicItem("extra blanket", nil, []string{"fancon"}, nil),
	plib.NewBasicItem("games", []string{"Con"}, nil),
	plib.NewBasicItem("sewing kit", []string{"Suiting", "BigTrip"}, nil),
	plib.NewBasicItem("grooming brush", []string{"Suiting"}, nil),
	plib.NewBasicItem("lint roller", []string{"Suiting"}, nil),
	plib.NewBasicItem("glue gun", []string{"Suiting"}, nil),
	plib.NewBasicItem("dryer sheets", []string{"Suiting"}, nil),
	plib.NewBasicItem("suiting outfit", []string{"Suiting"}, nil),
}

func init() {
	plib.RegisterItems("Unmentionables", unmentionables)
	plib.RegisterItems("Suiting", suiting)
	plib.RegisterItems("Playtime", playtime)
}
