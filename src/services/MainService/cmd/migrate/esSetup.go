package main

import (
	"fmt"
	"github.com/200-status-ok/main-backend/src/MainService/DBConfiguration"
	"github.com/200-status-ok/main-backend/src/MainService/Repository/ElasticSearch"
)

func main() {
	esClient := ElasticSearch.NewPosterES(DBConfiguration.GetElastic())
	err := esClient.DeletePosterIndex()
	if err != nil {
		fmt.Println(err)
	}
	err = esClient.CreatePosterIndex()

	if err != nil {
		fmt.Println(err)
	}
}
