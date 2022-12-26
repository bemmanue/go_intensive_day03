package main

import (
	"elastic/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Data struct {
	Name   string
	Places []db.Place
}

func main() {
	http.HandleFunc("/api/recommend", handler)
	log.Fatal(http.ListenAndServe("localhost:8888", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		writeError(w, "Error parsing query")
		return
	}

	lat, err := strconv.ParseFloat(values["lat"][0], 64)
	if err != nil {
		writeError(w, "Error parsing query")
		return
	}

	lon, err := strconv.ParseFloat(values["lon"][0], 64)
	if err != nil {
		writeError(w, "Error parsing query")
		return
	}

	places, err := db.GetPlaces(lat, lon)
	if err != nil {
		writeError(w, "error getting places")
		return
	}

	data := Data{
		Name:   "Recommendation",
		Places: places,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		writeError(w, "error getting places")
		return
	}
}

func writeError(w http.ResponseWriter, message string) {
	d := struct {
		Error string
	}{fmt.Sprint(message)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}
