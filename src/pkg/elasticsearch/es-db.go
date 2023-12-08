package elasticsearch

import (
	"fmt"
	elastic "github.com/elastic/go-elasticsearch/v8"
	"os"
)

var elasticClient *elastic.Client

// Connects to the Elasticsearch client using the provided connection string.
func connectElastic(connection string) (*elastic.Client, error) {
	client, err := elastic.NewClient(elastic.Config{
		Addresses: []string{
			connection,
		},
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func setupDatabase(appEnv string) {
	var esConnection string
	esConnection = "http://elasticsearch:9200/"
	if appEnv == "production" {
		esConnection = "http://elastic:jnL5TU7gFHucjbclKKd2AQnH@main-es:9200/"
	}
	fmt.Print("Connecting to Elasticsearch at ", esConnection, "\n")
	elasticClient, _ = connectElastic(esConnection)
}

func init() {
	appEnv := os.Getenv("APP_ENV2")
	setupDatabase(appEnv)
}

// GetElastic returns the Elasticsearch client instance.
func GetElastic() *elastic.Client {
	return elasticClient
}
