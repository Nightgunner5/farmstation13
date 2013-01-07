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
