package main

import (
	"fmt"
	"github.com/403-access-denied/main-backend/src/MainService/DBConfiguration"
	"github.com/403-access-denied/main-backend/src/MainService/Repository/ElasticSearch"
)

func main() {
	esClient := ElasticSearch.NewPosterES(DBConfiguration.GetElastic())
	err := esClient.DeleteIndex()
	if err != nil {
		fmt.Println(err)
	}
	err = esClient.CreateIndex()

	if err != nil {
		fmt.Println(err)
	}
}
