package main

import (
	"encoding/json"
	"ex02/db"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
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
	log.Fatal(http.ListenAndServe("server-cert:8888", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {

	if strings.HasPrefix(r.URL.RawQuery, "page=") {
		page, err := strconv.Atoi(strings.TrimPrefix(r.URL.RawQuery, "page="))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		places, total, err := db.GetPlaces(page)
		if err != nil {
			d := struct {
				Error string
			}{fmt.Sprintf("Invalid page value: %d", page)}
			err = json.NewEncoder(w).Encode(d)
			w.WriteHeader(http.StatusBadRequest)
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
			log.Fatalln(err)
		}
	}
}
