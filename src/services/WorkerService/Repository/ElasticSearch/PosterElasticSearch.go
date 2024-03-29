package ElasticSearch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/200-status-ok/main-backend/src/pkg/utils"
	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type PosterES struct {
	es *elastic.Client
}

func NewPosterES(es *elastic.Client) *PosterES {
	return &PosterES{es: es}
}

func (p *PosterES) UpdatePosterState(state string, docID int) error {
	updateJSON := fmt.Sprintf(`{
		"doc": {
			"state": "%s"
		}
	}`, state)

	updateReq := esapi.UpdateRequest{
		Index:      "posters",
		DocumentID: strconv.Itoa(docID),
		Body:       strings.NewReader(updateJSON),
		Refresh:    "true",
	}

	updateRes, err := updateReq.Do(context.Background(), p.es)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(updateRes.Body)

	return nil
}

func (p *PosterES) UpdateTags(result map[string]string) error {
	for k, v := range result {
		script := fmt.Sprintf("if (ctx._source.tags != null) { ctx._source.tags.stream().filter(tag -> tag.name == '%s').forEach(tag -> tag.state = '%s')}",
			k, v)
		body := map[string]interface{}{
			"script": map[string]interface{}{
				"source": script,
				"lang":   "painless",
			},
		}
		appEnv := os.Getenv("APP_ENV3")
		var url string
		if appEnv == "development" {
			url = utils.ReadFromEnvFile(".env", "LOCAL_ELASTIC_URL")
		} else if appEnv == "production" {
			url = utils.ReadFromEnvFile(".env", "PRODUCTION_ELASTIC_URL")
		}

		url = url + "posters/_update_by_query?refresh=true"

		jsonBody, err := json.Marshal(body)
		resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonBody)))
		if err != nil {
			return err
		}

		err = resp.Body.Close()
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			if resp.StatusCode == http.StatusConflict {
				for attempt := 0; attempt < 3; attempt++ {
					resp, err = http.Post(url, "application/json", strings.NewReader(string(jsonBody)))
					if err != nil {
						return err
					}
					err = resp.Body.Close()
					if err != nil {
						return err
					}
					if resp.StatusCode == http.StatusOK {
						break
					}
				}
			} else {
				return fmt.Errorf("unexpected status code %d", resp.StatusCode)
			}
		}
	}

	return nil
}
