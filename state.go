package main

import (
	"sync"
	"time"
)

var (
	state struct {
		Planters  []*Planter
		Harvested map[string]uint
		SeedTypes []string
		sync.RWMutex
	}
)

func init() {
	state.Harvested = make(map[string]uint)
	state.Harvested["Compost"] = 50
	for i := 0; i < 20; i++ {
		p := new(Planter)
		p.Defaults()
		p.Solution.Water = 200
		state.Planters = append(state.Planters, p)
	}
	for i := range Crops {
		if Crops[i].Type != Weed {
			state.SeedTypes = append(state.SeedTypes, Crops[i].Name)
		}
	}

	go func() {
		for {
			time.Sleep(time.Second)
			state.Lock()
			for _, p := range state.Planters {
				p.Tick()
			}
			state.Unlock()
		}
	}()
}
