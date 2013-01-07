package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		io.WriteString(w, Interface)
	})

	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/state" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		stateLock.RLock()
		defer stateLock.RUnlock()
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
	
</body>
</html>`
