package Migrate

import (
	"github.com/200-status-ok/main-backend/src/MainService/Cmd/DB"
	"github.com/200-status-ok/main-backend/src/MainService/Repository"
	"github.com/200-status-ok/main-backend/src/MainService/Repository/ElasticSearch"
	"github.com/200-status-ok/main-backend/src/pkg/elasticsearch"
)

func InsertAllPostersInES() {
	db, _ := DB.GetDB()
	repository := Repository.NewPosterRepository(db)
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
