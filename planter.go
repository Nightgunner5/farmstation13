package main

import (
	"math/rand"
)

type Solution struct {
	Water       float32 // decreases dehydration
	Compost     float32 // increases health and yield
	ToxicSlurry float32 // decreases health and yield, increases dehydration and mutation chance.

	// Magic plant formulas
	Mutriant float32 // increases mutation chance, decreases health
	GroBoost float32 // decreases growth time, increases dehydration
	TopCrop  float32 // increases yield, decreases number of harvests
}

type Planter struct {
	*Crop
	Solution
	Dehydration  float32 // if this reaches +-100, the plant dies. Negative values are drowning.
	Health       float32 // if this is 0, the plant is dead.
	YieldScale   float32 // multiplier for yield.
	TimeScale    float32 // multiplier for time (higher is slower).
	GrowthCycle  uint16  // if this is 0, the plant is harvestable.
	HarvestsLeft uint16  // ten times the estimated number of remaining harvests.
	Mutation     float32 // if this reaches 100, the plant mutates. high values cause damage to the plant.
}

func (p *Planter) Defaults() {
	p.Dehydration = 0
	p.Health = 100
	p.YieldScale = 1
	p.TimeScale = 1
	if p.Crop != nil {
		p.HarvestsLeft = p.Crop.Harvests
		p.GrowthCycle = uint16(p.Crop.Time)
	}
}

func moveTowards(current *float32, goal, scale float32) {
	if *current > goal {
		*current -= scale * rand.Float32()
		if *current < goal {
			*current = goal
		}
	} else {
		*current += scale * rand.Float32()
		if *current > goal {
			*current = goal
		}
	}
}

func (p *Planter) Tick() {
	if p.Crop == nil {
		if rand.Intn(100) == 0 {
			p.Crop = Weeds[rand.Intn(len(Weeds))]
			p.Defaults()
		}
		return
	}

	if p.Health <= 0 {
		return
	}

	if p.Mutation >= 100 {
		p.Health += 50
		p.Mutation = 0
		if len(p.Crop.Mutations) > 0 {
			p.Crop = &p.Crop.Mutations[rand.Intn(len(p.Crop.Mutations))]
			p.HarvestsLeft = p.Crop.Harvests
			moveTowards(&p.Solution.ToxicSlurry, 100, 10)
		}
	} else if p.Mutation > 50 {
		moveTowards(&p.Health, 0, 2.5)
	}

	if p.GrowthCycle != 0 {
		growth := uint16(rand.Intn(int(10 / p.TimeScale)))
		if p.GrowthCycle < growth {
			p.GrowthCycle = 0

			if p.Crop.Name == "Slurrypod" {
				p.Health = 0
				for _, p := range state.Planters {
					p.Solution.ToxicSlurry += rand.Float32()*10 + 10
				}
			} else if p.Crop.Name == "Pulsating Mass" {
				p.Health = 0
				for _, p := range state.Planters {
					p.Solution.Mutriant += 10000
				}
			}
		} else {
			p.GrowthCycle -= growth
		}
	} else if p.Crop.Name == "Radweed" {
		for _, p := range state.Planters {
			moveTowards(&p.Dehydration, 100, 10)
		}
	} else if p.Crop.Name == "Lasher" {
		for _, other := range state.Planters {
			if other != p {
				moveTowards(&other.Health, 0, 10)
			}
		}
	} else if p.Crop.Name == "Creeper" {
		for _, other := range state.Planters {
			if other.Crop == nil && rand.Intn(100) == 0 {
				other.Crop = p.Crop
				other.Health = p.Health
				other.Dehydration = p.Dehydration
				other.Defaults()
			}
		}
	}

	switch p.Crop.Time {
	default:
		fallthrough
	case VeryFast:
		moveTowards(&p.Solution.Water, 0, 30)
	case Fast:
		moveTowards(&p.Solution.Water, 0, 20)
	case Average:
		moveTowards(&p.Solution.Water, 0, 5)
	case Slow:
		moveTowards(&p.Solution.Water, 0, 2)
	case VerySlow:
		moveTowards(&p.Solution.Water, 0, 1)
	}
	moveTowards(&p.Solution.Compost, 0, 0.1)
	moveTowards(&p.Solution.ToxicSlurry, 0, 0.1)
	moveTowards(&p.Solution.Mutriant, 0, 0.1)
	moveTowards(&p.Solution.GroBoost, 0, 0.1)
	moveTowards(&p.Solution.TopCrop, 0, 0.1)
	moveTowards(&p.TimeScale, 1, 0.01)
	moveTowards(&p.YieldScale, 1, 0.01)
	moveTowards(&p.Mutation, 100, 0.1)

	if p.Solution.Water > 200 {
		moveTowards(&p.Dehydration, -100, 10)
	} else if p.Solution.Water <= 0 {
		moveTowards(&p.Dehydration, 100, 10)
	} else {
		moveTowards(&p.Dehydration, 0, 10)
	}

	if p.Solution.Compost > 0 {
		moveTowards(&p.Health, 100, 5)
		moveTowards(&p.TimeScale, 0, 0.01)
		if p.YieldScale < 10 { // don't move down if there's TopCrop
			moveTowards(&p.YieldScale, 10, 0.05)
		}
	}

	if p.Solution.ToxicSlurry > 0 {
		moveTowards(&p.Health, 0, 10)
		moveTowards(&p.Dehydration, 100, 10)
		moveTowards(&p.YieldScale, 0, 0.05)
		moveTowards(&p.Mutation, 100, 5)
	}

	if p.Solution.Mutriant > 0 || p.Solution.GroBoost > 0 || p.Solution.TopCrop > 0 {
		moveTowards(&p.Health, 100, 1)

		if p.Solution.Mutriant > 0 {
			moveTowards(&p.Health, 0, 10)
			moveTowards(&p.Mutation, 100, 25)
		}

		if p.Solution.GroBoost > 0 {
			moveTowards(&p.TimeScale, 0.001, 0.1)
			moveTowards(&p.Solution.Water, 0, 20) // extra gulp
		}

		if p.Solution.TopCrop > 0 {
			moveTowards(&p.YieldScale, 20, 0.05)
			if rand.Intn(100) == 0 && p.HarvestsLeft != 0 {
				p.HarvestsLeft--
			}
		}
	}

	if p.Dehydration < -50 || p.Dehydration > 50 {
		moveTowards(&p.Health, 0, 10)
	}
	if p.Dehydration <= -100 || p.Dehydration >= 100 {
		p.Health = 0
	}
	if p.Health > 100 {
		p.Health = 100
	}
}

func (p *Planter) Harvest() int {
	if p.Crop == nil {
		return -3
	}

	if p.Health <= 0 && (p.GrowthCycle != 0 || p.HarvestsLeft == 0 || p.Crop.Yield == 0) {
		p.Crop = nil
		p.Defaults()
		return -2
	}

	if p.GrowthCycle != 0 || p.Crop.Yield == 0 {
		return -1
	}

	yield := rand.Intn(int(float32(p.Crop.Yield)*p.YieldScale + 1))

	harvestsUsed := uint16(rand.Intn(19) + 1)
	if harvestsUsed > p.HarvestsLeft {
		p.HarvestsLeft = 0
		p.Health = 0
		return yield
	}

	p.HarvestsLeft -= harvestsUsed
	p.GrowthCycle = uint16(p.Crop.Time)
	return yield
}
