package main

import (
	"sync"
	"time"
)

var (
	state     = make(map[*Planter]bool)
	stateLock sync.RWMutex
)

func init() {
	for i := 0; i < 20; i++ {
		p := new(Planter)
		p.Defaults()
		p.Solution.Water = 200
		state[p] = true
	}

	go func() {
		for {
			time.Sleep(time.Second)
			stateLock.Lock()
			for p := range state {
				p.Tick()
			}
			stateLock.Unlock()
		}
	}()
}
