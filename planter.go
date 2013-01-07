package main

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
	Health       float32 // if this reaches 0, the plant dies.
	YieldScale   float32 // multiplier for yield.
	TimeScale    float32 // multiplier for time (higher is slower).
	GrowthCycle  uint16
	HarvestsLeft uint16 // ten times the estimated number of remaining harvests.

	// TODO: mutation
}
