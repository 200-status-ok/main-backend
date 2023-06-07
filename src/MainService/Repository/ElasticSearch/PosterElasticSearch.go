package ElasticSearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/403-access-denied/main-backend/src/MainService/DTO"
	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"io"
	"net/url"
	"strconv"
	"strings"
)

type ESPoster struct {
	es *elastic.Client
}

func NewPosterES(es *elastic.Client) *ESPoster {
	return &ESPoster{es: es}
}

func (p *ESPoster) CreatePosterIndex() error {
	req := esapi.IndicesCreateRequest{
		Index: "posters",
		Body: strings.NewReader(
			`{
				  "settings": {
					"analysis": {
					  "char_filter": {
						"zero_width_spaces": {
							"type":       "mapping",
							"mappings": [ "\\u200C=>\\u0020"]
						}
					  },
					  "filter": {
						"persian_stop": {
						  "type":       "stop",
						  "stopwords":  "_persian_"
						}
					  },
					  "analyzer": {
						"rebuilt_persian": {
						  "tokenizer":     "standard",
						  "char_filter": [ "zero_width_spaces" ],
						  "filter": [
							"lowercase",
							"decimal_digit",
							"arabic_normalization",
							"persian_normalization",
							"persian_stop"
						  ]
						}
					  }
					}
				  },
				"mappings": {
					"properties": {
					  "poster": {
						"properties": {
						  "title": {
							"type": "search_as_you_type",
							"analyzer": "persian"
						  },
						  "description": {
							"type": "search_as_you_type",
							"analyzer": "persian"
						  },
						  "status": {
							"type": "keyword"
						  },
						  "tel_id": {
							"type": "keyword"
						  },
						  "user_phone": {
							"type": "keyword"
						  },
						  "alert": {
							"type": "boolean"
						  },
						  "chat": {
							"type": "boolean"
						  },
						  "award": {
							"type": "float"
						  },
						  "user_id": {
							"type": "integer"
						  },
						  "state": {
							"type": "keyword"
						  },
						  "special_type": {
							"type": "keyword"
						  },
						  "created_at": {
							"type": "date"
						  },
						  "updated_at": {
							"type": "date"
						  }
						}
					  },
					  "addresses": {
						"properties": {
						  "province": {
							"type": "text",
							"analyzer": "persian"
						  },
						  "city": {
							"type": "text",
							"analyzer": "persian"
						  },
						  "address_detail": {
							"type": "text",
							"analyzer": "persian"
						  },
						  "location": {
                            "type": "geo_point"
                          }
						}
					  },
					  "tags": {
						"properties": {
						  "id": {
							"type": "integer"
						  },
						  "name": {
							"type": "text",
							"analyzer": "persian"
						  },
						  "state": {
							"type": "keyword"
						  }
						}
					  },
                      "images": {
						"type": "keyword"
					  }}
				}
            }
		`),
	}
	res, err := req.Do(context.Background(), p.es)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	if err != nil {
		return err
	}
	return nil
}

func (p *ESPoster) DeletePosterIndex() error {
	req := esapi.IndicesDeleteRequest{
		Index: []string{"posters"},
	}

	res, err := req.Do(context.Background(), p.es)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	if err != nil {
		return err
	}
	return nil
}

func (p *ESPoster) InsertPoster(poster *DTO.ESPosterDTO) error {
	jsonPoster, err := json.Marshal(poster)
	if err != nil {
		return err
	}
	req := esapi.IndexRequest{
		Index:      "posters",
		DocumentID: strconv.Itoa(int(poster.ID)),
		Body:       strings.NewReader(string(jsonPoster)),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), p.es)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	if err != nil {
		return err
	}
	return nil
}

func (p *ESPoster) InsertAllPosters(posters []*DTO.ESPosterDTO) error {
	var bulkRequest []string
	for _, poster := range posters {
		jsonPoster, err := json.Marshal(poster)
		if err != nil {
			return err
		}
		bulkRequest = append(bulkRequest, fmt.Sprintf(`{ "index" : { "_index" : "posters", "_id" : "%d" } }`, poster.ID))
		bulkRequest = append(bulkRequest, string(jsonPoster))
	}

	req := esapi.BulkRequest{
		Index:   "posters",
		Body:    strings.NewReader(strings.Join(bulkRequest, "\n") + "\n"),
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), p.es)
	fmt.Println(res.StatusCode)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(res.Body)

	if err != nil {
		return err
	}

	return nil
}

func (p *ESPoster) UpdatePoster(poster map[string]interface{}, posterID int) error {
	jsonPoster, err := json.Marshal(poster)
	if err != nil {
		return err
	}

	req := esapi.UpdateRequest{
		Index:      "posters",
		DocumentID: strconv.Itoa(posterID),
		Body:       strings.NewReader(string(jsonPoster)),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), p.es)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	if err != nil {
		return err
	}
	return nil
}

func (p *ESPoster) DeletePoster(posterID int) error {
	req := esapi.DeleteRequest{
		Index:      "posters",
		DocumentID: strconv.Itoa(posterID),
	}

	res, err := req.Do(context.Background(), p.es)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(res.Body)

	if err != nil {
		return err
	}
	return nil
}

func (p *ESPoster) DeletePosterByUserID(userID int) error {
	deleteQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"user_id": map[string]interface{}{
					"value": userID,
				},
			},
		},
	}
	jsonData, err := json.Marshal(deleteQuery)
	if err != nil {
		return err
	}

	req := esapi.DeleteByQueryRequest{
		Index: []string{"posters"},
		Body:  strings.NewReader(string(jsonData)),
	}

	res, err := req.Do(context.Background(), p.es)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	return nil
}

type SearchHits struct {
	Hits Hits `json:"hits"`
}

type Hits struct {
	Total Total  `json:"total"`
	Hits  []*Hit `json:"hits"`
}

type Total struct {
	Value int `json:"value"`
}

type Hit struct {
	Source *DTO.ESPosterDTO `json:"_source"`
}

func (p *ESPoster) GetPosters(filterObject DTO.FilterObject) ([]*DTO.ESPosterDTO, int, error) {
	getPostersFields := make(map[string]interface{})
	getPostersFields["track_scores"] = true
	getPostersFields["from"] = filterObject.Offset
	getPostersFields["size"] = filterObject.PageSize
	getPostersFields["sort"] = []map[string]interface{}{
		{
			"_score": map[string]interface{}{
				"order": "desc",
			},
		},
		{
			filterObject.SortBy: map[string]interface{}{
				"order": filterObject.Sort,
			},
		},
	}
	getPostersFields["query"] = map[string]interface{}{
		"bool": map[string]interface{}{
			"should": []map[string]interface{}{
				{
					"term": map[string]interface{}{
						"special_type": "premium",
					},
				},
			},
			"must": []map[string]interface{}{
				{
					"terms": map[string]interface{}{
						"special_type": []string{"premium", "normal"},
					},
				},
				{
					"terms": map[string]interface{}{
						"state": []string{"rejected", "accepted", "pending"},
					},
				},
				{
					"terms": map[string]interface{}{
						"status": []string{"lost", "found"},
					},
				},
			},
		},
	}
	if filterObject.Status != "both" {
		statusQuery := map[string]interface{}{
			"term": map[string]interface{}{
				"status": filterObject.Status,
			},
		}
		getPostersFields["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{})[2] =
			statusQuery
	}
	if filterObject.State != "all" {
		stateQuery := map[string]interface{}{
			"term": map[string]interface{}{
				"state": filterObject.State,
			},
		}
		getPostersFields["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{})[1] =
			stateQuery

	}
	if filterObject.SpecialType != "all" {
		specialQuery := map[string]interface{}{
			"term": map[string]interface{}{
				"special_type": filterObject.SpecialType,
			},
		}
		getPostersFields["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{})[0] =
			specialQuery
	}
	if filterObject.OnlyAwards {
		awardQuery := map[string]interface{}{
			"range": map[string]interface{}{
				"award": map[string]interface{}{
					"gt": 0,
				},
			},
		}
		getPostersFields["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(
			getPostersFields["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
			awardQuery)
	}
	if filterObject.Lat != 0 && filterObject.Lon != 0 {
		geoQuery := map[string]interface{}{
			"geo_distance": map[string]interface{}{
				"distance": "10km",
				"addresses.location": map[string]interface{}{
					"lat": filterObject.Lat,
					"lon": filterObject.Lon,
				},
			},
		}
		getPostersFields["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(
			getPostersFields["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
			geoQuery)
	}
	if filterObject.SearchPhrase != "" {
		decodedSearchPhrase, _ := url.QueryUnescape(filterObject.SearchPhrase)
		searchQuery := map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": decodedSearchPhrase,
				"type":  "bool_prefix",
				"fields": []string{"title", "description", "title._2gram",
					"description._2gram", "title._3gram", "description._3gram"},
				"fuzziness": "AUTO",
			},
		}
		getPostersFields["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(
			getPostersFields["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
			searchQuery)
	}

	jsonData, err := json.Marshal(getPostersFields)
	if err != nil {
		return []*DTO.ESPosterDTO{}, 0, err
	}

	res, err := p.es.Search(
		p.es.Search.WithIndex("posters"),
		p.es.Search.WithBody(strings.NewReader(string(jsonData))),
		p.es.Search.WithTrackScores(true),
	)
	if err != nil {
		return []*DTO.ESPosterDTO{}, 0, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(res.Body)

	if res.IsError() {
		if res.StatusCode == 404 || res.StatusCode == 400 {
			return []*DTO.ESPosterDTO{}, 0, nil
		}
		return []*DTO.ESPosterDTO{}, 0, errors.New(fmt.Sprintf("Search request failed: %s", res.Status()))
	}

	var result SearchHits
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return []*DTO.ESPosterDTO{}, 0, err
	}

	postersHits := result.Hits.Hits
	var getPosters []*DTO.ESPosterDTO

	for _, v := range postersHits {
		getPosters = append(getPosters, v.Source)
	}

	return getPosters, result.Hits.Total.Value, nil
}
