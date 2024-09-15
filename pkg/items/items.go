package items

import "github.com/ywwg/packingdb/pkg/packinglib"

func RegisterAllItems(r packinglib.Registry) {
	r.RegisterItems("Art", art)
	r.RegisterItems("Business", business)
	r.RegisterItems("Camping", campStuff)
	r.RegisterItems("Clothing", clothing)
	r.RegisterItems("Documents", documents)
	r.RegisterItems("Dog", dog)
	r.RegisterItems("Electrical", electrical)
	r.RegisterItems("Entertainment", entertainment)
	r.RegisterItems("Flight Stuff", flightSupplies)
	r.RegisterItems("Food", food)
	r.RegisterItems("Performing", performing)

	// Sport
	r.RegisterItems("Bicycling", bicycling)
	r.RegisterItems("Climbing", climbing)
	r.RegisterItems("Skiing", skiing)
	r.RegisterItems("Boat", boating)

	// Swim
	r.RegisterItems("Swim", waterStuff)
	r.RegisterItems("Tubing", tubing)

	r.RegisterItems("To Buy", toBuy)
	r.RegisterItems("Tasks", tasks)
	r.RegisterItems("Toiletries", toiletries)
	r.RegisterItems("Tools", tools)

	r.RegisterItems("Unmentionables", unmentionables)
	r.RegisterItems("Suiting", suiting)
	r.RegisterItems("Playtime", playtime)
}
