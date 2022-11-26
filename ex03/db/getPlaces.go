package db

import (
	"bytes"
	"errors"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"strconv"
	"time"
)

type Place struct {
	Name    string
	Address string
	Phone   string
}

func GetPlaces(page int) ([]Place, int, error) {
	var (
		places   []Place
		total    int
		scrollID string
	)

	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating client: %s", err)
	}

	res, err := es.Search(
		es.Search.WithIndex("places"),
		es.Search.WithSize(10),
		es.Search.WithScroll(time.Minute),
	)
	if err != nil {
		log.Fatalf("Error searching document: %s", err)
	}

	myjson := read(res.Body)
	res.Body.Close()

	scrollID = gjson.Get(myjson, "_scroll_id").String()
	total, _ = strconv.Atoi(gjson.Get(myjson, "hits.total.value").String())

	if total <= (page-1)*10 {
		return nil, total, errors.New("the page doesn't exist")
	}

	for i := 1; i < page; i++ {

		res, err := es.Scroll(es.Scroll.WithScrollID(scrollID), es.Scroll.WithScroll(time.Minute))
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		if res.IsError() {
			log.Fatalf("Error response: %s", res)
		}

		myjson = read(res.Body)
		res.Body.Close()
	}

	hits := gjson.Get(myjson, "hits.hits").Array()
	for _, hit := range hits {
		var place Place
		place.Name = hit.Map()["_source"].Map()["Name"].Map()["Type"].String()
		place.Address = hit.Map()["_source"].Map()["Address"].Map()["Type"].String()
		place.Phone = hit.Map()["_source"].Map()["Phone"].Map()["Type"].String()
		places = append(places, place)
	}

	return places, total, nil
}

func read(r io.Reader) string {
	var b bytes.Buffer
	b.ReadFrom(r)
	return b.String()
}
