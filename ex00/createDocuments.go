package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/olivere/elastic/v7"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

type Type struct {
	Type string `json:"type"`
}

type GeoType struct {
	Type elastic.GeoPoint `json:"type"`
}

type Doc struct {
	Name     Type    `json:"name"`
	Address  Type    `json:"address"`
	Phone    Type    `json:"phone"`
	Location GeoType `json:"location"`
}

var documentsCount int = 25

func createDocuments() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating client: %s", err)
	}

	for i := 1; i <= documentsCount; i++ {
		var mydoc Doc
		lat := rand.Float64() + 55
		lon := rand.Float64() + 37

		mydoc.Name = Type{"Restaurant_№" + strconv.Itoa(i)}
		mydoc.Address = Type{"City_№" + strconv.Itoa(i)}
		mydoc.Phone = Type{"Phone_№" + strconv.Itoa(i)}
		mydoc.Location = GeoType{elastic.GeoPoint{Lat: lat, Lon: lon}}

		myjson := jsonStruct(mydoc)

		request := esapi.IndexRequest{
			Index:      "places",
			DocumentID: strconv.Itoa(i),
			Body:       strings.NewReader(myjson),
			Refresh:    "true",
		}

		response, err := request.Do(context.Background(), es)
		if err != nil {
			log.Fatalf("Error creating request: %s", err)
		}
		defer response.Body.Close()

		if response.IsError() {
			log.Fatalln("Error indexing document")
		}

		var res map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&res); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}

		fmt.Println("Status:", response.Status())
	}
}

func jsonStruct(doc Doc) string {
	res, err := json.Marshal(&doc)
	if err != nil {
		log.Fatalf("Error marshalling structure: %s", err)
	}
	return string(res)
}
