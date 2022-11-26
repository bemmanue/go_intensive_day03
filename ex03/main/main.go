package main

import (
	"ex03/db"
	"log"
	"net/http"
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
	http.HandleFunc("/api/recommend", handler)
	log.Fatal(http.ListenAndServe("localhost:8888", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println(r)
}
