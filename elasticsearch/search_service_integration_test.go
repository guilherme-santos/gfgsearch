// +build integration

package elasticsearch_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/guilherme-santos/gfgsearch"
	"github.com/guilherme-santos/gfgsearch/elasticsearch"

	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	ctx := context.Background()
	esClient, cleanup := esClient(t)
	defer cleanup()

	searchSvc := elasticsearch.NewSearchService(esClient)
	searchSvc.InitMapping(ctx)

	loadESData(t, esClient, "products.json")

	resp, err := searchSvc.Search(ctx, "", gfgsearch.Options{
		Page:    1,
		PerPage: 5,
	})

	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.EqualValues(t, 10, resp.Total)
	assert.Len(t, resp.Data, 5)

	for i, p := range resp.Data {
		id := i + 1
		assert.Equal(t, fmt.Sprintf("product %d", id), p.Title)
		assert.Equal(t, fmt.Sprintf("brand %d", id), p.Brand)
		assert.EqualValues(t, id*100, p.Price)
		assert.EqualValues(t, id, p.Stock)
	}
}
