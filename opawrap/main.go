package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"

	"opawrap/queryeval"
)

// TODO: sistemare gli err

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /query", handleQuery)

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleQuery(w http.ResponseWriter, r *http.Request) {

	baseURL := &url.URL{
		Scheme: "http",
		Host:   "localhost:8081",
		Path:   "data",
	}

	// retrieve input to send to the policy engine
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var input any
	err = json.Unmarshal(body, &input)
	if err != nil {
		http.Error(w, "invalid input of the query", http.StatusBadRequest)
		return
	}

	var data map[string]any
	getState(baseURL, &data, w)

	state, result := queryeval.OPA(data, input, w, r.Context())

	updateState(baseURL, state, w)

	// Return only necessary output (i.e. without state part) to user
	output, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

/*
 * Makes a GET request to datastore's API, in order to retrieve the actual state
 *
 * baseURL - the URL of datastore server
 * state - variable that will contain the actual state
 * w - to handle http errors
 */
func getState(baseURL *url.URL, state *map[string]any, w http.ResponseWriter) {
	client := http.Client{}
	resp, err := client.Get(baseURL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dataBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//var data map[string]any
	err = json.Unmarshal(dataBytes, &state)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

/*
 * Makes a POST request to datastore's API, in order to update the state
 *
 * baseURL - the URL of datastore server
 * state - variable that contains the new state
 * w - to handle http errors
 */
func updateState(baseURL *url.URL, state map[string]any, w http.ResponseWriter) {
	for key, value := range state {
		baseURL.Path = path.Join("/data", key)
		bodyBytes, err := json.Marshal(value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		body := io.NopCloser(bytes.NewBuffer(bodyBytes))
		client := http.Client{}
		req, err := http.NewRequest(http.MethodPut, baseURL.String(), body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = client.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
