package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/botany/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/botany/" {
			http.NotFound(w, r)
			return
		}

		io.WriteString(w, Interface)
	})

	log.Fatal(http.ListenAndServe(":26301", nil))
}

const Interface = `<!DOCTYPE html>
<html>
<head>
	<title>Farm Station 13</title>
	<style>
body {
	padding-bottom: 100px;
}
#harvested {
	position: absolute;
	right: 8px;
	top: 8px;
}
#seeds {
	position: fixed;
	bottom: 0;
	background-color: rgba(255, 255, 255, 0.5);
}
button {
	border: 1px solid #777;
	border-radius: 5px;
	background: -moz-linear-gradient(top, #eee, #ccc, #bbb, #999);
	background: -webkit-linear-gradient(top, #eee, #ccc, #bbb, #999);
	background: linear-gradient(top, #eee, #ccc, #bbb, #999);
}
button:active {
	box-shadow: 0 0 5px #000 inset;
}
[id^="planter-"]>button:first-child {
	border-radius: 5px 0 0 5px;
	margin-right: 0;
}
[id^="planter-"]>button:first-child+button,
[id^="planter-"]>button:first-child+button+button {
	border-radius: 0;
	border-left: 0;
	margin-left: 0;
	margin-right: 0;
}
[id^="planter-"]>button:first-child+button+button+button {
	border-radius: 0 5px 5px 0;
	border-left: 0;
	margin-left: 0;
}
.plantName {
	display: inline-block;
	width: 100px;
}
.no-water {
	background-color: #f77;
	padding: 3px;
	border-radius: 5px;
	box-shadow: 0 0 5px 5px #fff inset;
}
.low-water {
	background-color: #ff7;
	padding: 3px;
	border-radius: 5px;
	box-shadow: 0 0 5px 5px #fff inset;
}
.good-water {
	background-color: #7f7;
	padding: 3px;
	border-radius: 5px;
	box-shadow: 0 0 5px 5px #fff inset;
}
.drown-water {
	background-color: #77f;
	padding: 3px;
	border-radius: 5px;
	box-shadow: 0 0 5px 5px #fff inset;
}
	</style>
</head>
<body>
	<div id="harvested"></div>
	<script>
function handle(msg) {
	var state = JSON.parse(msg.data);

	for (var i = 0; i < state.Planters.length; ++i) {
		var p = state.Planters[i];
		if (p.Name) {
			planter(i, p.Name, p.Health, p);
		} else {
			planter(i, '', -1, p);
		}
	}

	// Seed Types
	var seeds = document.getElementById('seeds');
	if (seeds == null) {
		seeds = document.createElement('div');
		seeds.id = 'seeds';
		state.SeedTypes.forEach(function(s) {
			var seed = document.createElement('button');
			seed.innerText = s;
			seed.onclick = function() {
				ws.send(JSON.stringify({Action:'Plant', Crop:s}));
			};
			seeds.appendChild(seed);
		});
		document.body.appendChild(seeds);
	}

	// Harvested
	for (var crop in state.Harvested) {
		var amount = state.Harvested[crop];
		var el = document.getElementById('harvested-' + crop);
		if (amount > 0) {
			if (el == null) {
				el = document.createElement('div');
				el.id = 'harvested-' + crop;
				if (crop != 'Compost') {
					var mulch = document.createElement('button');
					mulch.innerText = 'Mulch';
					mulch.onclick = (function(crop) {
						return function() {
							ws.send(JSON.stringify({Action:'Mulch', Crop:crop}));
						};
					})(crop);
					el.appendChild(mulch);
					el.appendChild(document.createTextNode(' '));
				}
				el.appendChild(document.createElement('strong'));
				el.appendChild(document.createTextNode(' ' + crop));

				document.getElementById('harvested').appendChild(el);
			}
			el.querySelector('strong').innerText = amount + '×';
		} else {
			if (el != null) {
				el.parentNode.removeChild(el);
			}
		}
	}
}

function planter(i, name, health, data) {
	var p = document.getElementById('planter-' + i);
	if (p == null) {
		p = document.createElement('div');
		p.id = 'planter-' + i;
		var drain = document.createElement('button');
		drain.onclick = function() {
			ws.send(JSON.stringify({Action:'Drain', Planter:i}));
		};
		drain.innerText = 'D';
		drain.title = 'Drain';
		p.appendChild(drain);

		var chainsaw = document.createElement('button');
		chainsaw.onclick = function() {
			ws.send(JSON.stringify({Action:'Chainsaw', Planter:i}));
		};
		chainsaw.innerText = 'X';
		chainsaw.title = 'Chainsaw';
		p.appendChild(chainsaw);

		var water = document.createElement('button');
		water.onclick = function() {
			ws.send(JSON.stringify({Action:'Water', Planter:i}));
		};
		water.innerText = 'W';
		water.title = 'Water';
		p.appendChild(water);

		var compost = document.createElement('button');
		compost.onclick = function() {
			ws.send(JSON.stringify({Action:'Compost', Planter:i}));
		};
		compost.innerText = 'C';
		compost.title = 'Compost';
		p.appendChild(compost);

		var plantName = document.createElement('strong');
		plantName.className = 'plantName';
		p.appendChild(plantName);

		p.appendChild(document.createTextNode(' - '));
		p.appendChild(document.createElement('span'));
		p.appendChild(document.createTextNode(' '));
		p.appendChild(document.createElement('em'));
		document.body.appendChild(p);
	}

	p.querySelector('strong').innerText = name;

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

	if (contents.length == 0) contents.push('Nothing');

	var solutionMeter = p.querySelector('span');
	solutionMeter.innerText = solution + ' units of ' + contents.join(', ');
	solutionMeter.className = data.Water > 200 ? 'drown-water' :
		data.Water > 75 ? 'good-water' :
		data.Water > 0 ? 'low-water' : 'no-water';

	var status = p.querySelector('em');
	status.innerText = '';
	if (health > 50) {
		if (data.Dehydration > 50) {
			status.appendChild(document.createTextNode('(Unhealthy - Dehydrated)'));
		} else if (data.Dehydration < -50) {
			status.appendChild(document.createTextNode('(Unhealthy - Drowning)'));
		} else {
			status.appendChild(document.createTextNode('(Healthy)'));
		}
	} else if (health > 0) {
		if (data.Dehydration > 50) {
			status.appendChild(document.createTextNode('(Unhealthy - Dehydrating)'));
		} else if (data.Dehydration < -50) {
			status.appendChild(document.createTextNode('(Unhealthy - Drowning)'));
		} else {
			status.appendChild(document.createTextNode('(Unhealthy)'));
		}
	} else if (health == 0) {
		status.appendChild(document.createTextNode('(Dead '));

		var clear = document.createElement('button');
		clear.onclick = function() {
			ws.send(JSON.stringify({Action:'Harvest', Planter:i}));	
		};
		clear.innerText = data.GrowthCycle == 0 && data.HarvestsLeft != 0 && data.Yield != 0 ? 'Harvest' : 'Clear';
		status.appendChild(clear);

		status.appendChild(document.createTextNode(')'));
	}

	if (health > 0) {
		if (data.GrowthCycle > data.Time / 2) {
			status.appendChild(document.createTextNode(' (Sprouting)'));
		} else if (data.GrowthCycle == 0) {
			if (data.Yield != 0) {
				status.appendChild(document.createTextNode(' (Mature '));

				var harvest = document.createElement('button');
				harvest.onclick = function() {
					ws.send(JSON.stringify({Action:'Harvest', Planter:i}));	
				};
				harvest.innerText = 'Harvest';
				status.appendChild(harvest);

				status.appendChild(document.createTextNode(')'));
			}
		} else {
			status.appendChild(document.createTextNode(' (Growing)'));
		}
	}
}

var ws = new WebSocket('ws://' + location.host + '/botany/sock');

ws.onclose = function closed() {
	ws = new WebSocket('ws://' + location.host + '/botany/sock');
	ws.onmessage = handle;
	ws.onclose = closed;
};
ws.onmessage = handle;
	</script>
</body>
</html>`
