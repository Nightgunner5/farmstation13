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
	Dehydration  float32 // if this reaches 100, the plant dies.
	Health       float32 // if this is 0, the plant is dead.
	YieldScale   float32 // multiplier for yield.
	TimeScale    float32 // multiplier for time (higher is slower).
	GrowthCycle  uint16
	HarvestsLeft uint16 // ten times the estimated number of remaining harvests.

	// TODO: mutation
}

func (p *Planter) Defaults() {
	p.Dehydration = 0
	p.Health = 100
	p.YieldScale = 1
	p.TimeScale = 1
	p.GrowthCycle = 0
	p.HarvestsLeft = 0
}

func (p *Planter) Tick() {
	// TODO
}

func (p *Planter) Harvest() int {
	if p.Crop == nil {
		return -3
	}

	if p.GrowthCycle != 0 || p.Crop.Yield == 0 {
		return -2
	}

	if p.Health <= 0 {
		p.Crop = nil
		p.Defaults()
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
	p.GrowthCycle = uint16(float32(p.Crop.Time) * p.TimeScale)
	return yield
}
