package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

var filename = "data.json"

func main() {
	store := make(map[string]any)
	err := loadFromFile(&store)
	if err != nil {
		log.Fatal()
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /data", handleData(store))
	mux.HandleFunc("PUT /data/{key}", handleUpdate(store))

	log.Fatal(http.ListenAndServe(":8081", mux))
}

func loadFromFile(store *map[string]any) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &store)
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

		err = saveToFile(store)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func saveToFile(store map[string]any) error {
	data, err := json.Marshal(store)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
