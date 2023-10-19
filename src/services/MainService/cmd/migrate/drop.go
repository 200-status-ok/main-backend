package main

import (
	Model2 "github.com/200-status-ok/main-backend/src/MainService/Model"
	"github.com/200-status-ok/main-backend/src/pkg/pgsql"
)

func main() {
	var models []interface{}
	models = append(models, &Model2.User{})
	models = append(models, &Model2.Poster{})
	models = append(models, &Model2.Tag{})
	models = append(models, &Model2.Conversation{})
	models = append(models, &Model2.Message{})
	models = append(models, &Model2.Image{})
	models = append(models, &Model2.Address{})
	models = append(models, &Model2.MarkedPoster{})
	models = append(models, &Model2.PosterReport{})

	err := pgsql.DropModel(models)
	if err != nil {
		return
	}
}
