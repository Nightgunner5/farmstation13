package main

import (
	"sync"
	"time"
)

var (
	state struct {
		Planters  []*Planter
		Harvested map[string]uint
		sync.RWMutex
	}
)

func init() {
	state.Harvested = make(map[string]uint)
	for i := 0; i < 20; i++ {
		p := new(Planter)
		p.Defaults()
		p.Solution.Water = 200
		state.Planters = append(state.Planters, p)
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
