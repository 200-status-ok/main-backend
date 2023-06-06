package ElasticSearch

import (
	"context"
	"fmt"
	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"io"
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
