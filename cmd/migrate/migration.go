package main

import (
	"github.com/403-access-denied/main-backend/src/DBConfiguration"
	"github.com/403-access-denied/main-backend/src/Model"
)

func main() {
	Migration()
}

func Migration() {
	var models []interface{}
	// add the model to migrate the database
	models = append(models, &Model.User{})
	models = append(models, &Model.Poster{})
	models = append(models, &Model.Category{})
	models = append(models, &Model.Conversation{})
	models = append(models, &Model.Message{})
	models = append(models, &Model.Image{})
	models = append(models, &Model.Address{})
	models = append(models, &Model.MarkedPoster{})
	DBConfiguration.InitDB()
	DBConfiguration.MigrateModel(models)
	DBConfiguration.CloseDB()
}
