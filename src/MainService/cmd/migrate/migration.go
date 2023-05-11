package main

import (
	"github.com/403-access-denied/main-backend/src/MainService/DBConfiguration"
	Model2 "github.com/403-access-denied/main-backend/src/MainService/Model"
)

func main() {
	var models []interface{}
	// add the model to migrate the database
	models = append(models, &Model2.User{})
	models = append(models, &Model2.Poster{})
	models = append(models, &Model2.Tag{})
	models = append(models, &Model2.Conversation{})
	models = append(models, &Model2.Message{})
	models = append(models, &Model2.Image{})
	models = append(models, &Model2.Address{})
	models = append(models, &Model2.MarkedPoster{})
	models = append(models, &Model2.PosterReport{})
	models = append(models, &Model2.Payment{})
	models = append(models, &Model2.Admin{})
	//DBConfiguration.init()
	DBConfiguration.MigrateModel(models)
	//DBConfiguration.CloseDB()
}
