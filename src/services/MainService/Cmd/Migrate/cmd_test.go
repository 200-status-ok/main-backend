package Migrate

import (
	"github.com/200-status-ok/main-backend/src/MainService/Model"
	"github.com/200-status-ok/main-backend/src/MainService/Repository/ElasticSearch"
	"github.com/200-status-ok/main-backend/src/pkg/elasticsearch"
	"github.com/200-status-ok/main-backend/src/pkg/pgsql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModelsMigrate(t *testing.T) {
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
	assert.NoError(t, err)

	assert.True(t, pgsql.GetDB().Migrator().HasTable(&Model.User{}))
	assert.True(t, pgsql.GetDB().Migrator().HasTable(&Model.Poster{}))
	assert.True(t, pgsql.GetDB().Migrator().HasTable(&Model.Tag{}))
	assert.True(t, pgsql.GetDB().Migrator().HasTable(&Model.Conversation{}))
	assert.True(t, pgsql.GetDB().Migrator().HasTable(&Model.Message{}))
	assert.True(t, pgsql.GetDB().Migrator().HasTable(&Model.Image{}))
	assert.True(t, pgsql.GetDB().Migrator().HasTable(&Model.Address{}))
	assert.True(t, pgsql.GetDB().Migrator().HasTable(&Model.MarkedPoster{}))
	assert.True(t, pgsql.GetDB().Migrator().HasTable(&Model.PosterReport{}))
	assert.True(t, pgsql.GetDB().Migrator().HasTable(&Model.Payment{}))
	assert.True(t, pgsql.GetDB().Migrator().HasTable(&Model.Admin{}))
}

func TestESSetup(t *testing.T) {
	esClient := ElasticSearch.NewPosterES(elasticsearch.GetElastic())
	err := esClient.DeletePosterIndex()
	assert.NoError(t, err)
	err = esClient.CreatePosterIndex()
	assert.NoError(t, err)
}
