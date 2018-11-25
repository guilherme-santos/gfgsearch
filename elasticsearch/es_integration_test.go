// +build integration

package elasticsearch_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/guilherme-santos/gfgsearch/elasticsearch"

	"github.com/stretchr/testify/assert"
)

func esClient(t *testing.T) (*elasticsearch.Client, func()) {
	esURL := os.Getenv("ELASTICSEARCH_URL")

	esClient, err := elasticsearch.NewClient(esURL)
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
