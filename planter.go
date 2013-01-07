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

	// TODO: mutation
}

func (p *Planter) Defaults() {
	p.Dehydration = 0
	p.Health = 100
	p.YieldScale = 1
	p.TimeScale = 1
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
			p.GrowthCycle = uint16(p.Crop.Time)
		}
		return
	}

	if p.GrowthCycle != 0 {
		growth := uint16(rand.Intn(int(10 / p.TimeScale)))
		if p.GrowthCycle < growth {
			p.GrowthCycle = 0

			if p.Crop.Name == "Slurrypod" {
				p.Crop = nil
				p.Defaults()
				go func() {
					stateLock.Lock()
					defer stateLock.Unlock()
					for _, p := range state {
						p.Solution.ToxicSlurry += rand.Float32()*10 + 10
					}
				}()
			}
		} else {
			p.GrowthCycle -= growth
		}
	} else if p.Crop.Name == "Radweed" {
		go func() {
			stateLock.Lock()
			defer stateLock.Unlock()
			for _, p := range state {
				moveTowards(&p.Dehydration, 100, 10)
			}
		}()
	} else if p.Crop.Name == "Creeper" {
		go func() {
			stateLock.Lock()
			defer stateLock.Unlock()
			c := p.Crop
			for _, p := range state {
				if p.Crop == nil && rand.Intn(100) == 0 {
					p.Crop = c
					p.Defaults()
				}
			}
		}()
	}

	if p.Health <= 0 {
		return
	}

	moveTowards(&p.Solution.Water, 0, 5)
	moveTowards(&p.Solution.Compost, 0, 5)
	moveTowards(&p.Solution.ToxicSlurry, 0, 5)
	moveTowards(&p.Solution.Mutriant, 0, 5)
	moveTowards(&p.Solution.GroBoost, 0, 5)
	moveTowards(&p.Solution.TopCrop, 0, 5)
	moveTowards(&p.TimeScale, 1, 0.01)
	moveTowards(&p.YieldScale, 1, 0.01)

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
		// TODO: mutation
	}

	if p.Solution.Mutriant > 0 || p.Solution.GroBoost > 0 || p.Solution.TopCrop > 0 {
		moveTowards(&p.Health, 100, 1)

		if p.Solution.Mutriant > 0 {
			moveTowards(&p.Health, 0, 10)
			// TODO: mutation
		}

		if p.Solution.GroBoost > 0 {
			moveTowards(&p.TimeScale, 0.001, 0.01)
			moveTowards(&p.Solution.Water, 0, 10) // extra gulp
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
