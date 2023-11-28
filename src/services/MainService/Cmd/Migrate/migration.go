package Migrate

import (
	"github.com/200-status-ok/main-backend/src/MainService/Model"
	"github.com/200-status-ok/main-backend/src/pkg/pgsql"
)

func ModelsMigrate() {
	var models []interface{}

	models = append(models, &Model.User{})
	models = append(models, &Model.Poster{})
	models = append(models, &Model.Tag{})
	models = append(models, &Model.Conversation{})
	models = append(models, &Model.Message{})
	models = append(models, &Model.Image{})
	models = append(models, &Model.Address{})
	models = append(models, &Model.MarkedPoster{})
	models = append(models, &Model.PosterReport{})
	models = append(models, &Model.Payment{})
	models = append(models, &Model.Admin{})

	err := pgsql.MigrateModel(models)
	if err != nil {
		return
	}

}
