package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		io.WriteString(w, Interface)
	})

	http.HandleFunc("/harvest/", func(w http.ResponseWriter, r *http.Request) {
		i, err := strconv.ParseInt(r.URL.Path[len("/harvest/"):], 10, 64)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		w.WriteHeader(http.StatusNoContent)

		state.Lock()
		defer state.Unlock()
		if 0 <= int(i) && int(i) < len(state.Planters) {
			crop := state.Planters[int(i)].Crop
			amount := state.Planters[int(i)].Harvest()
			if amount < 0 {
				// TODO
				return
			}
			state.Harvested[crop.Name] += uint(amount)
		}
	})

	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/state" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		state.RLock()
		defer state.RUnlock()
		json.NewEncoder(w).Encode(state)
	})

	log.Fatal(http.ListenAndServe(":26301", nil))
}

const Interface = `<!DOCTYPE html>
<html>
<head>
	<title>Farm Station 13</title>
</head>
<body>
	<script>
function planter(i, name, health, data) {
	var p = document.createElement('div');
	var n = document.createElement('strong');
	n.innerText = name;
	p.appendChild(n);

	n = document.createTextNode(' - ');
	p.appendChild(n);

	var solution = 0, contents = [];
	solution += data.Water;
	if (data.Water > 0) contents.push('Water');
	solution += data.Compost;
	if (data.Compost > 0) contents.push('Compost');
	solution += data.ToxicSlurry;
	if (data.ToxicSlurry > 0) contents.push('Toxic Slurry');
	solution += data.Mutriant;
	if (data.Mutriant > 0) contents.push('Mutriant');
	solution += data.GroBoost;
	if (data.GroBoost > 0) contents.push('Gro-Boost');
	solution += data.TopCrop;
	if (data.TopCrop > 0) contents.push('TopCrop');
	solution = Math.round(solution * 100) / 100;

	n = document.createTextNode(solution + ' units of ' + contents.join(', '));
	p.appendChild(n);

	if (health > 50) {
		n = document.createTextNode(' (Healthy)');
		p.appendChild(n);
	} else if (health > 0) {
		n = document.createTextNode(' (Unhealthy)');
		p.appendChild(n);
	} else if (health == 0) {
		n = document.createTextNode(' (');
		p.appendChild(n);

		n = document.createElement('a');
		n.href = '/harvest/' + i;
		n.innerText = 'Dead';
		p.appendChild(n);

		n = document.createTextNode(')');
		p.appendChild(n);
	}

	if (health > 0) {
		if (data.GrowthCycle > data.Time / 2) {
			n = document.createTextNode(' (Sprouting)');
			p.appendChild(n);
		} else if (data.GrowthCycle == 0) {
			if (data.Yield != 0) {
				n = document.createTextNode(' (');
				p.appendChild(n);

				n = document.createElement('a');
				n.href = '/harvest/' + i;
				n.innerText = 'Harvestable';
				p.appendChild(n);

				n = document.createTextNode(')');
				p.appendChild(n);
			}
		} else {
			n = document.createTextNode(' (Growing)');
			p.appendChild(n);
		}
	}

	document.body.appendChild(p);
}

function harvested(crop, amount) {
	var harvested = document.getElementById('harvested');
	if (!harvested) {
		harvested = document.createElement('div');
		harvested.style.position = 'absolute';
		harvested.style.top = '8px';
		harvested.style.right = '8px';
		document.body.appendChild(harvested);
	}

	var h = document.createElement('div');
	var n = document.createElement('strong');
	n.innerText = amount + 'x';
	h.appendChild(n);

	n = document.createTextNode(' ' + crop);
	h.appendChild(n);

	harvested.appendChild(h);
}

setInterval(function() {
	var xhr = new XMLHttpRequest();
	xhr.open('get', '/state', true);
	xhr.onload = function() {
		var state = JSON.parse(xhr.responseText);
		document.body.innerHTML = '';
		state.Planters.forEach(function(p, i) {
			if (p.Name) {
				planter(i, p.Name, p.Health, p);
			} else {
				planter(i, 'Empty', -1, p);
			}
		});
		for (var crop in state.Harvested) {
			var amount = state.Harvested[crop];
			if (amount > 0) harvested(crop, amount);
		}
	};
	xhr.send();
}, 1000);
	</script>
</body>
</html>`
