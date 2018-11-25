package elasticsearch

import (
	"net/http"
	"time"

	"github.com/olivere/elastic"
)

// Timeout sets the timeout when communicate to ElasticSearch
var Timeout = 1 * time.Second

type Client struct {
	*elastic.Client
}

// NewClient create a new elasticsearch client.
func NewClient(url string) (*Client, error) {
	httpClient := &http.Client{
		// In case elasticsearch (server side) doesn't answer in the
		// timeout specified, we're creating a client timeout (our side),
		// just to be safe that our connection don't be hold.
		Timeout: 2 * Timeout,
	}

	esClient, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
		elastic.SetHttpClient(httpClient),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client: esClient,
	}, nil
}
