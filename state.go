package main

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"math/rand"
	"net/http"
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

	stateNotify = sync.NewCond(new(sync.Mutex))
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

			stateNotify.L.Lock()
			stateNotify.Broadcast()
			stateNotify.L.Unlock()
		}
	}()

	http.Handle("/botany/sock", websocket.Handler(HandleSocket))
}

func HandleSocket(ws *websocket.Conn) {
	addr := ws.Request().RemoteAddr

	go func() {
		// Write
		for {
			stateNotify.L.Lock()
			stateNotify.Wait()
			stateNotify.L.Unlock()

			state.RLock()
			err := websocket.JSON.Send(ws, state)
			state.RUnlock()

			if err != nil {
				return
			}
		}
	}()

	// Read
	for {
		var packet struct {
			Action  string
			Crop    string
			Planter int
		}

		err := websocket.JSON.Receive(ws, &packet)

		if err != nil {
			log.Println(addr, err)

			ws.Close()
			return
		}

		switch packet.Action {
		case "Harvest":
			state.Lock()
			if packet.Planter >= 0 && packet.Planter < len(state.Planters) {
				p := state.Planters[packet.Planter]
				crop := p.Crop
				amount := p.Harvest()
				if amount < 0 {
					// TODO
				} else {
					log.Println(addr, "harvested", amount, crop.Name)
					state.Harvested[crop.Name] += uint(amount)
				}
			}
			state.Unlock()

			stateNotify.L.Lock()
			stateNotify.Broadcast()
			stateNotify.L.Unlock()

		case "Drain":
			state.Lock()
			if packet.Planter >= 0 && packet.Planter < len(state.Planters) {
				p := state.Planters[packet.Planter]
				crop := p.Crop
				cropName := "empty"
				if crop != nil {
					cropName = crop.Name
				}
				log.Printf("%s drained planter %d %+v (crop: %s)", addr, packet.Planter, p.Solution, cropName)
				p.Solution = Solution{}
			}
			state.Unlock()

			stateNotify.L.Lock()
			stateNotify.Broadcast()
			stateNotify.L.Unlock()

		case "Chainsaw":
			state.Lock()
			if packet.Planter >= 0 && packet.Planter < len(state.Planters) {
				p := state.Planters[packet.Planter]
				if p.Crop != nil {
					log.Println(addr, "chainsawed planter", packet.Planter, "(crop:", p.Crop.Name+")")
					if p.Health <= 0 {
						p.Crop = nil
					} else {
						moveTowards(&p.Health, 0, 100)
					}
				}
			}
			state.Unlock()

			stateNotify.L.Lock()
			stateNotify.Broadcast()
			stateNotify.L.Unlock()

		case "Water":
			state.Lock()
			if packet.Planter >= 0 && packet.Planter < len(state.Planters) {
				p := state.Planters[packet.Planter]
				p.Solution.Water += 60
				crop := p.Crop
				cropName := "empty"
				if crop != nil {
					cropName = crop.Name
				}
				log.Println(addr, "watered planter", packet.Planter, "(new level:", p.Solution.Water, "units)", "(crop:", cropName+")")
			}
			state.Unlock()

			stateNotify.L.Lock()
			stateNotify.Broadcast()
			stateNotify.L.Unlock()

		case "Compost":
			state.Lock()
			if packet.Planter >= 0 && packet.Planter < len(state.Planters) {
				p := state.Planters[packet.Planter]
				compost := state.Harvested["Compost"]
				if compost >= 10 {
					compost = 10
					switch rand.Intn(3) {
					case 0:
						p.Solution.Mutriant += rand.Float32() * 50
					case 1:
						p.Solution.GroBoost += rand.Float32() * 50
					case 2:
						p.Solution.TopCrop += rand.Float32() * 50
					}
				}
				log.Println(addr, "added compost to planter", packet.Planter)
				state.Harvested["Compost"] -= compost
				p.Solution.Compost += float32(compost)
			}
			state.Unlock()

			stateNotify.L.Lock()
			stateNotify.Broadcast()
			stateNotify.L.Unlock()

		case "Mulch":
			state.Lock()
			if h, ok := state.Harvested[packet.Crop]; ok && h > 0 {
				var amount uint = 1
				if h > 50 {
					amount = 10
				}
				if h > 250 {
					amount = 100
				}
				log.Println(addr, "mulched", amount, packet.Crop)
				state.Harvested[packet.Crop] -= amount
				state.Harvested["Compost"] += amount
			}
			state.Unlock()

			stateNotify.L.Lock()
			stateNotify.Broadcast()
			stateNotify.L.Unlock()

		case "Plant":
			for i := range Crops {
				if Crops[i].Type == Weed || Crops[i].Name != packet.Crop {
					continue
				}

				state.Lock()
				for _, p := range state.Planters {
					if p.Crop != nil {
						continue
					}
					log.Println(addr, "planted", packet.Crop)
					p.Crop = &Crops[i]
					p.Defaults()
					break
				}
				state.Unlock()

				stateNotify.L.Lock()
				stateNotify.Broadcast()
				stateNotify.L.Unlock()
				break
			}
		}
	}
}
