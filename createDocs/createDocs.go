package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"log"
	"strconv"
	"strings"
)

type Type struct {
	Type string
}

type Doc struct {
	Name    Type
	Address Type
	Phone   Type
}

func jsonStruct(doc Doc) string {
	res, err := json.Marshal(&doc)
	if err != nil {
		log.Fatalf("Error marshalling structure: %s", err)
	}
	return string(res)
}

func main() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating client: %s", err)
	}

	for i := 1; i <= 15; i++ {
		var mydoc Doc
		mydoc.Name = Type{"Person_№" + strconv.Itoa(i)}
		mydoc.Address = Type{"City_№" + strconv.Itoa(i)}
		mydoc.Phone = Type{"Phone_№" + strconv.Itoa(i)}

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
		} else {
			var res map[string]interface{}
			if err := json.NewDecoder(response.Body).Decode(&res); err != nil {
				log.Printf("Error parsing the response body: %s", err)
			} else {
				fmt.Println("Status:", response.Status())
			}
		}
	}
}
