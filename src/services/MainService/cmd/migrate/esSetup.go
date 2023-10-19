package main

import (
	"fmt"
	"github.com/200-status-ok/main-backend/src/MainService/Repository/ElasticSearch"
	"github.com/200-status-ok/main-backend/src/pkg/elasticsearch"
)

func main() {
	esClient := ElasticSearch.NewPosterES(elasticsearch.GetElastic())
	err := esClient.DeletePosterIndex()
	if err != nil {
		fmt.Println(err)
	}
	err = esClient.CreatePosterIndex()

	if err != nil {
		fmt.Println(err)
	}
}
