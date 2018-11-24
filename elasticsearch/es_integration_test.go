// +build integration

package elasticsearch_test

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/guilherme-santos/gfgsearch/elasticsearch"

	"github.com/olivere/elastic"
	"github.com/stretchr/testify/assert"
)

func esClient(t *testing.T) (*elastic.Client, func()) {
	esURL := os.Getenv("ELASTICSEARCH_URL")
	httpClient := &http.Client{
		Timeout: 2 * time.Second,
	}

	esClient, err := elastic.NewClient(
		elastic.SetURL(esURL),
		elastic.SetSniff(false),
		elastic.SetHttpClient(httpClient),
		// elastic.SetTraceLog(log.New(os.Stdout, "", 0)),
	)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	// save current es index
	curESIndex := elasticsearch.Index
	// Update rand seed
	rand.Seed(time.Now().UnixNano())
	// generate temporary es index
	elasticsearch.Index = curESIndex + "_" + fmt.Sprint(rand.Int63())
	// This func should be called to restore the index name and also delete
	// the temporary index created.
	cleanupFn := func() {
		esClient.DeleteIndex(elasticsearch.Index).Do(context.Background())
		elasticsearch.Index = curESIndex
	}

	return esClient, cleanupFn
}
