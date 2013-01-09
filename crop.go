package main

type Type uint8

const (
	Fruit Type = iota
	Vegetable
	Ingredient
	Herb
	Weed
)

type Time uint16

const (
	VeryFast Time = 1 << (iota + 4)
	Fast
	Average
	Slow
	VerySlow
)

type Crop struct {
	Name     string
	Type     Type
	Yield    uint16 // 0 means unharvestable
	Harvests uint16 // 10 times the average number of harvests. 1 for definite single harvest.
	Time     Time

	Mutations []Crop
}

const infinite = ^uint16(0)

var Crops = []Crop{
	{
		Name:     "Apple",
		Type:     Fruit,
		Yield:    5,
		Harvests: 80,
		Time:     VerySlow,

		Mutations: []Crop{
			{
				Name:     "Snapple",
				Type:     Fruit,
				Yield:    2,
				Harvests: 120,
				Time:     Average,
			},
		},
	},
	{
		Name:     "Asomna",
		Type:     Herb,
		Yield:    3,
		Harvests: 60,
		Time:     Average,
	},
	{
		Name:     "Avocado",
		Type:     Vegetable,
		Yield:    2,
		Harvests: 20,
		Time:     Slow,
	},
	{
		Name:     "Banana",
		Type:     Fruit,
		Yield:    5,
		Harvests: 50,
		Time:     Slow,
	},
	{
		Name:     "Carrot",
		Type:     Vegetable,
		Yield:    3,
		Harvests: 30,
		Time:     Average,

		Mutations: []Crop{
			{
				Name:     "Apathyroot",
				Type:     Vegetable,
				Yield:    1,
				Harvests: 160,
				Time:     VerySlow,
			},
		},
	},
	{
		Name:     "Chili",
		Type:     Ingredient,
		Yield:    3,
		Harvests: 40,
		Time:     Slow,

		Mutations: []Crop{
			{
				Name:     "Chilly",
				Type:     Ingredient,
				Yield:    20,
				Harvests: 10,
				Time:     VerySlow,
			},
		},
	},
	{
		Name:     "Commol",
		Type:     Herb,
		Yield:    3,
		Harvests: 10,
		Time:     Average,
	},
	{
		Name:     "Contusine",
		Type:     Herb,
		Yield:    3,
		Harvests: 10,
		Time:     Average,
	},
	{
		Name:     "Corn",
		Type:     Vegetable,
		Yield:    5,
		Harvests: 30,
		Time:     Fast,
	},
	{
		Name:     "Creeper",
		Type:     Weed,
		Yield:    0,
		Harvests: 0,
		Time:     Average,
	},
	{
		Name:     "Eggplant",
		Type:     Fruit,
		Yield:    4,
		Harvests: 30,
		Time:     Average,
	},
	{
		Name:     "Garlic",
		Type:     Ingredient,
		Yield:    3,
		Harvests: 70,
		Time:     Average,
	},
	{
		Name:     "Grapes",
		Type:     Fruit,
		Yield:    5,
		Harvests: 80,
		Time:     Slow,
	},
	{
		Name:     "Lasher",
		Type:     Weed,
		Yield:    0,
		Harvests: 0,
		Time:     Slow,
	},
	{
		Name:     "Lemon",
		Type:     Fruit,
		Yield:    4,
		Harvests: 50,
		Time:     VerySlow,

		Mutations: []Crop{
			{
				Name:     "Pink Lemon",
				Type:     Fruit,
				Yield:    5,
				Harvests: 40,
				Time:     Slow,
			},
		},
	},
	{
		Name:     "Lettuce",
		Type:     Vegetable,
		Yield:    7,
		Harvests: 1,
		Time:     Fast,
	},
	{
		Name:     "Lime",
		Type:     Fruit,
		Yield:    4,
		Harvests: 50,
		Time:     VerySlow,
	},
	{
		Name:     "Melon",
		Type:     Fruit,
		Yield:    1,
		Harvests: 25,
		Time:     Slow,

		Mutations: []Crop{
			{
				Name:     "George Melons",
				Type:     Fruit,
				Yield:    7,
				Harvests: 1,
				Time:     VerySlow,
			},
		},
	},
	{
		Name:     "Nureous",
		Type:     Herb,
		Yield:    3,
		Harvests: 30,
		Time:     Average,
	},
	{
		Name:     "Onion",
		Type:     Vegetable,
		Yield:    6,
		Harvests: 40,
		Time:     Average,
	},
	{
		Name:     "Orange",
		Type:     Fruit,
		Yield:    4,
		Harvests: 60,
		Time:     VerySlow,

		Mutations: []Crop{
			{
				Name:     "Blood Orange",
				Type:     Fruit,
				Yield:    2,
				Harvests: 500,
				Time:     VeryFast,
			},
		},
	},
	{
		Name:     "Peanut",
		Type:     Ingredient,
		Yield:    6,
		Harvests: 40,
		Time:     Average,
	},
	{
		Name:     "Potato",
		Type:     Vegetable,
		Yield:    5,
		Harvests: 1,
		Time:     Average,
	},
	{
		Name:     "Pumpkin",
		Type:     Vegetable,
		Yield:    2,
		Harvests: 10,
		Time:     Slow,

		Mutations: []Crop{
			{
				Name:     "Great Pumpkin",
				Type:     Vegetable,
				Yield:    1,
				Harvests: 1,
				Time:     VerySlow,
			},
		},
	},
	{
		Name:     "Radweed",
		Type:     Weed,
		Yield:    0,
		Harvests: 0,
		Time:     Slow,
	},
	{
		Name:     "Slurrypod",
		Type:     Weed,
		Yield:    0,
		Harvests: 0,
		Time:     Average,

		Mutations: []Crop{
			{
				Name:     "Pulsating Mass",
				Type:     Weed,
				Yield:    0,
				Harvests: 0,
				Time:     Fast,
			},
		},
	},
	{
		Name:     "Soy",
		Type:     Ingredient,
		Yield:    4,
		Harvests: 1, // social commentary lol
		Time:     Average,

		Mutations: []Crop{
			{
				Name:     "Soylent Tofu",
				Type:     Ingredient,
				Yield:    5,
				Harvests: 20,
				Time:     Slow,
			},
		},
	},
	{
		Name:     "Space Fungus",
		Type:     Weed,
		Yield:    3,
		Harvests: infinite,
		Time:     Slow,

		Mutations: []Crop{
			{
				Name:     "Space Shroom",
				Type:     Weed,
				Yield:    2,
				Harvests: infinite,
				Time:     VerySlow,
			},
			{
				Name:     "Space Mold",
				Type:     Weed,
				Yield:    4,
				Harvests: infinite,
				Time:     VerySlow,
			},
		},
	},
	{
		Name:     "Space Grass",
		Type:     Weed,
		Yield:    4,
		Harvests: infinite,
		Time:     VeryFast,

		Mutations: []Crop{
			{
				Name:     "Crab Grass",
				Type:     Weed,
				Yield:    1,
				Harvests: infinite,
				Time:     Fast,
			},
		},
	},
	{
		Name:     "Sugar Cane",
		Type:     Ingredient,
		Yield:    6,
		Harvests: 20,
		Time:     Fast,

		Mutations: []Crop{
			{
				Name:     "High Fructose Sugar-Like Sweetener Substance",
				Type:     Ingredient,
				Yield:    6,
				Harvests: 40,
				Time:     Average,
			},
		},
	},
	{
		Name:     "Synthmeat",
		Type:     Ingredient,
		Yield:    4,
		Harvests: 40,
		Time:     Fast,

		Mutations: []Crop{
			{
				Name:     "Synthbrain",
				Type:     Ingredient,
				Yield:    10,
				Harvests: 60,
				Time:     Slow,
			},
			{
				Name:     "Synthflesh",
				Type:     Ingredient,
				Yield:    1,
				Harvests: infinite,
				Time:     VerySlow,
			},
		},
	},
	{
		Name:     "Tomato",
		Type:     Fruit,
		Yield:    5,
		Harvests: 40,
		Time:     Slow,

		Mutations: []Crop{
			{
				Name:     "Thomas-ato",
				Type:     Fruit,
				Yield:    20,
				Harvests: 20,
				Time:     VerySlow,
			},
		},
	},
	{
		Name:     "Venne",
		Type:     Herb,
		Yield:    3,
		Harvests: 30,
		Time:     Average,
	},
	{
		Name:     "Wheat",
		Type:     Ingredient,
		Yield:    8,
		Harvests: 10,
		Time:     Fast,

		Mutations: []Crop{
			{
				Name:     "Steelwheat",
				Type:     Ingredient,
				Yield:    3,
				Harvests: 400,
				Time:     Average,
			},
		},
	},
}

var Weeds []*Crop

func init() {
	for i := range Crops {
		if Crops[i].Type == Weed {
			Weeds = append(Weeds, &Crops[i])
		}
	}
}
