package main

import (
	"encoding/json"
	"ex03/db"
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
		log.Println(err)
	}

	lat, err := strconv.ParseFloat(values["lat"][0], 64)
	if err != nil {
		log.Fatalln(err)
	}

	lon, err := strconv.ParseFloat(values["lon"][0], 64)
	if err != nil {
		log.Fatalln(err)
	}

	places, err := db.GetPlaces(lat, lon)
	if err != nil {
		log.Println(err)
	}

	data := Data{
		Name:   "Recommendation",
		Places: places,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Fatalln(err)
	}
}
