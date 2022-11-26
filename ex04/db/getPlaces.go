package db

import (
	"bytes"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/tidwall/gjson"
	"io"
	"log"
)

type Geo struct {
	Lat float64
	Lon float64
}

type Place struct {
	Name     string
	Address  string
	Phone    string
	Location Geo
}

func GetPlaces(lat, lon float64) ([]Place, error) {
	var places []Place
	var buf bytes.Buffer

	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating client: %s", err)
	}

	sort := map[string]interface{}{
		"sort": map[string]interface{}{
			"_geo_distance": map[string]interface{}{
				"location": map[string]interface{}{
					"lat": lat,
					"lon": lon,
				},
				"order":           "asc",
				"unit":            "km",
				"mode":            "min",
				"distance_type":   "arc",
				"ignore_unmapped": true,
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(sort); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	search, err := es.Search(
		es.Search.WithIndex("places"),
		es.Search.WithSize(3),
		es.Search.WithBody(&buf),
	)
	if err != nil {
		log.Fatalf("Error searching document: %s", err)
	}

	myjson := read(search.Body)
	search.Body.Close()

	hits := gjson.Get(myjson, "hits.hits").Array()
	for _, hit := range hits {
		var place Place
		place.Name = hit.Map()["_source"].Map()["name"].Map()["type"].String()
		place.Address = hit.Map()["_source"].Map()["address"].Map()["type"].String()
		place.Phone = hit.Map()["_source"].Map()["phone"].Map()["type"].String()
		place.Location.Lat = hit.Map()["_source"].Map()["location"].Map()["type"].Map()["lat"].Float()
		place.Location.Lon = hit.Map()["_source"].Map()["location"].Map()["type"].Map()["lon"].Float()
		places = append(places, place)
	}

	return places, nil
}

func read(r io.Reader) string {
	var b bytes.Buffer
	b.ReadFrom(r)
	return b.String()
}
