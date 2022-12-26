package main

import (
	"elastic/db"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
)

type Data struct {
	Title    string
	Total    int
	Places   []db.Place
	Previous int
	Next     int
	Last     int
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8888", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query().Get("page")
	page, err := strconv.Atoi(param)
	if err != nil {
		writeInvalidParamError(w, r.URL.String())
		return
	}

	places, total, err := db.GetPlaces(page)
	if err != nil {
		writeInvalidPageError(w, page)
		return
	}

	data := Data{
		Title:    "Pages",
		Total:    total,
		Places:   places,
		Previous: page - 1,
		Next:     page + 1,
		Last:     int(math.Ceil(float64(total) / 10.0)),
	}
	if page == data.Last {
		data.Next = 0
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		writeInternalError(w, page)
	}
}

func writeInvalidParamError(w http.ResponseWriter, param string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	d := struct {
		Error string
	}{fmt.Sprintf("Invalid parameter: %d", param)}
	json.NewEncoder(w).Encode(d)
}

func writeInvalidPageError(w http.ResponseWriter, page int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	d := struct {
		Error string
	}{fmt.Sprintf("Invalid page value: %d", page)}
	json.NewEncoder(w).Encode(d)
}

func writeInternalError(w http.ResponseWriter, page int) {
	w.WriteHeader(http.StatusBadRequest)
	d := struct {
		Error string
	}{fmt.Sprintf("Page %d: internal error", page)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}
