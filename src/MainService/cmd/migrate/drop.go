package main

import (
	"github.com/403-access-denied/main-backend/src/MainService/DBConfiguration"
	Model2 "github.com/403-access-denied/main-backend/src/MainService/Model"
)

func main() {
	var models []interface{}
	models = append(models, &Model2.User{})
	models = append(models, &Model2.Poster{})
	models = append(models, &Model2.Category{})
	models = append(models, &Model2.Conversation{})
	models = append(models, &Model2.Message{})
	models = append(models, &Model2.Image{})
	models = append(models, &Model2.Address{})
	models = append(models, &Model2.MarkedPoster{})
	models = append(models, &Model2.ChatRoom{})
	models = append(models, &Model2.PosterReport{})
	DBConfiguration.InitDB()
	DBConfiguration.DropModel(models)
	DBConfiguration.CloseDB()
}
