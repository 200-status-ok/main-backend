package main

import (
	"github.com/200-status-ok/main-backend/src/MainService/Repository"
	"github.com/200-status-ok/main-backend/src/MainService/Repository/ElasticSearch"
	"github.com/200-status-ok/main-backend/src/pkg/elasticsearch"
	"github.com/200-status-ok/main-backend/src/pkg/pgsql"
)

func main() {
	repository := Repository.NewPosterRepository(pgsql.GetDB())
	allESPosters, err := repository.GetAllESPosters()

	if err != nil {
		panic(err)
	}
	esPosterCli := ElasticSearch.NewPosterES(elasticsearch.GetElastic())
	err = esPosterCli.InsertAllPosters(allESPosters)
	if err != nil {
		panic(err)
	}
}
