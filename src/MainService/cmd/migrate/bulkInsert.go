package main

import (
	"github.com/403-access-denied/main-backend/src/MainService/DBConfiguration"
	"github.com/403-access-denied/main-backend/src/MainService/Repository"
	"github.com/403-access-denied/main-backend/src/MainService/Repository/ElasticSearch"
)

func main() {
	repository := Repository.NewPosterRepository(DBConfiguration.GetDB())
	allESPosters, err := repository.GetAllESPosters()

	if err != nil {
		panic(err)
	}
	esPosterCli := ElasticSearch.NewPosterES(DBConfiguration.GetElastic())
	err = esPosterCli.InsertAllPosters(allESPosters)
	if err != nil {
		panic(err)
	}
}
