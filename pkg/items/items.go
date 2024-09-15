package items

import "github.com/ywwg/packingdb/pkg/packinglib"

func AllItems() []packinglib.ItemList {
	return []packinglib.ItemList{
		{Name: "Art", Items: art},
		{Name: "Business", Items: business},
		{Name: "Camping", Items: campStuff},
		{Name: "Clothing", Items: clothing},
		{Name: "Documents", Items: documents},
		{Name: "Dog", Items: dog},
		{Name: "Electrical", Items: electrical},
		{Name: "Entertainment", Items: entertainment},
		{Name: "Flight Stuff", Items: flightSupplies},
		{Name: "Food", Items: food},
		{Name: "Performing", Items: performing},

		// Sport
		{Name: "Bicycling", Items: bicycling},
		{Name: "Climbing", Items: climbing},
		{Name: "Skiing", Items: skiing},
		{Name: "Boat", Items: boating},

		// Swim
		{Name: "Swim", Items: waterStuff},
		{Name: "Tubing", Items: tubing},

		{Name: "To Buy", Items: toBuy},
		{Name: "Tasks", Items: tasks},
		{Name: "Toiletries", Items: toiletries},
		{Name: "Tools", Items: tools},

		{Name: "Unmentionables", Items: unmentionables},
		{Name: "Suiting", Items: suiting},
		{Name: "Playtime", Items: playtime},
	}
}
