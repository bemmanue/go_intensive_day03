package main

import (
	"context"
	"github.com/olivere/elastic/v7"
	"log"
)

var mapping string = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings": {
		"properties": {
			"name": {
				"type": "text"
			},
			"address": {
				"type": "text"
			},
			"phone": {
				"type": "text"
			},
			"location": {
				"type": "geo_point"
			}
		}
	}
}
`

func main() {
	index := "places"

	// If elasticsearch works on host
	//es, err := elastic.NewClient(elastic.SetURL("server-cert:9200"))

	// In case of elasticsearch works on docker container
	es, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		log.Fatalf("Error creating the client: %s\n", err)
	}

	exists, err := es.IndexExists(index).Do(context.Background())
	if err != nil {
		log.Fatalf("Error checking index existence: %s\n", err)
	}

	if !exists {
		createIndex, err := es.CreateIndex(index).BodyString(mapping).Do(context.Background())
		if err != nil {
			log.Println("Error creating index")
		}
		if !createIndex.Acknowledged {
			log.Println("Adding index is not acknowledged")
		} else {
			log.Println("successfully created index")
		}
	} else {
		log.Println("Index already exist")
	}

	createDocuments()
}
