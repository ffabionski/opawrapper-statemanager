package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// TODO: rendere lo store persistente

func main() {
	store := make(map[string]any)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /data", handleData(store))
	mux.HandleFunc("PUT /data/{key}", handleUpdate(store))

	log.Fatal(http.ListenAndServe(":8081", mux))
}

func handleData(store map[string]any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(store)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func handleUpdate(store map[string]any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var data any
		err = json.Unmarshal(body, &data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		store[key] = data
		log.Printf("Set %q to %v", key, data)
		w.WriteHeader(http.StatusOK)
	}
}
