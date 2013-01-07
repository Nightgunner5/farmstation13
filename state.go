package main

import (
	"sync"
	"time"
)

var (
	state     []*Planter
	stateLock sync.RWMutex
)

func init() {
	for i := 0; i < 20; i++ {
		p := new(Planter)
		p.Defaults()
		p.Solution.Water = 200
		state = append(state, p)
	}

	go func() {
		for {
			time.Sleep(time.Second)
			stateLock.Lock()
			for _, p := range state {
				p.Tick()
			}
			stateLock.Unlock()
		}
	}()
}
