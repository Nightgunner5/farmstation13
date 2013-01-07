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
	VeryFast Time = 1 << (iota + 5)
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
}

const infinite = ^uint16(0)

var Crops = []Crop{
	{"Apple", Fruit, 5, 80, VerySlow},
	{"Asomna", Herb, 3, 60, Average},
	{"Avocado", Vegetable, 2, 20, Slow},
	{"Banana", Fruit, 5, 50, Slow},
	{"Carrot", Vegetable, 3, 30, Average},
	{"Chili", Ingredient, 3, 40, Slow},
	{"Commol", Herb, 3, 10, Average},
	{"Contusine", Herb, 3, 10, Average},
	{"Corn", Vegetable, 5, 30, Fast},
	{"Creeper", Weed, 0, 0, Average},
	{"Eggplant", Fruit, 4, 30, Average},
	{"Garlic", Ingredient, 3, 70, Average},
	{"Grape", Fruit, 5, 80, Slow},
	{"Lasher", Weed, 0, 0, Slow},
	{"Lemon", Fruit, 4, 50, VerySlow},
	{"Lettuce", Vegetable, 7, 1, Fast},
	{"Lime", Fruit, 4, 50, VerySlow},
	{"Maneater", Weed, 0, 0, Slow},
	{"Melon", Fruit, 1, 25, Slow},
	{"Nureous", Herb, 3, 30, Average},
	{"Onion", Vegetable, 6, 40, Average},
	{"Orange", Fruit, 4, 60, VerySlow},
	{"Peanut", Ingredient, 6, 40, Average},
	{"Potato", Vegetable, 5, 1, Average},
	{"Pumpkin", Vegetable, 2, 1, Slow},
	{"Radweed", Weed, 0, 0, Slow},
	{"Slurrypod", Weed, 0, 0, Average},
	{"Soy", Ingredient, 4, 1 /* social commentary lol */, Average},
	{"Space Fungus", Weed, 3, infinite, VeryFast},
	{"Space Grass", Weed, 4, infinite, Slow},
	{"Sugar Cane", Ingredient, 6, 20, Fast},
	{"Synthmeat", Ingredient, 4, 40, Fast},
	{"Tomato", Fruit, 5, 40, Slow},
	{"Venne", Herb, 3, 30, Average},
	{"Wheat", Ingredient, 8, 10, Fast},
}

var Weeds []*Crop

func init() {
	for i := range Crops {
		if Crops[i].Type == Weed {
			Weeds = append(Weeds, &Crops[i])
		}
	}
}
