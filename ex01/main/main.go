package main

import (
	"elastic/db"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

type Data struct {
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
		Total:    total,
		Places:   places,
		Previous: page - 1,
		Next:     page + 1,
		Last:     int(math.Ceil(float64(total) / 10.0)),
	}
	if page == data.Last {
		data.Next = 0
	}

	tmpl, err := template.ParseFiles("./../templates/index.html")
	if err != nil {
		writeInternalError(w, page)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		writeInternalError(w, page)
		return
	}
}

func writeInvalidParamError(w http.ResponseWriter, param string) {
	w.WriteHeader(http.StatusBadRequest)
	tmpl, _ := template.ParseFiles("./../templates/invalid_param.html")
	tmpl.Execute(w, struct{ Param string }{Param: param})
}

func writeInvalidPageError(w http.ResponseWriter, page int) {
	w.WriteHeader(http.StatusBadRequest)
	tmpl, _ := template.ParseFiles("./../templates/invalid_param.html")
	tmpl.Execute(w, struct{ Page int }{Page: page})
}

func writeInternalError(w http.ResponseWriter, page int) {
	w.WriteHeader(http.StatusBadRequest)
	tmpl, _ := template.ParseFiles("./../templates/internal_error.html")
	tmpl.Execute(w, struct{ Page int }{Page: page})
}
