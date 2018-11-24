// +build integration

package elasticsearch_test

import (
	"context"
	"path/filepath"
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

	err := searchSvc.LoadFile(filepath.Join("testdata", "products.json"))
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	resp, err := searchSvc.Search(ctx, "", gfgsearch.Options{
		Page:    1,
		PerPage: 5,
	})

	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.EqualValues(t, 10, resp.Total)
	assert.Len(t, resp.Data, 5)
}

func TestSearch_FilterByProductTitle(t *testing.T) {
	ctx := context.Background()
	esClient, cleanup := esClient(t)
	defer cleanup()

	searchSvc := elasticsearch.NewSearchService(esClient)
	searchSvc.InitMapping(ctx)

	err := searchSvc.LoadFile(filepath.Join("testdata", "products.json"))
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	resp, err := searchSvc.Search(ctx, "shirt", gfgsearch.Options{
		Page:    1,
		PerPage: 5,
	})

	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.EqualValues(t, 3, resp.Total)
	assert.Len(t, resp.Data, 3)

	products := map[string]gfgsearch.Product{
		"basic t-shirt": {
			Title: "basic t-shirt",
			Brand: "hackett london",
			Price: 100,
			Stock: 10,
		},
		"printed t-shirt": {
			Title: "printed t-shirt",
			Brand: "tommy hilfiger",
			Price: 200,
			Stock: 9,
		},
		"business shirt": {
			Title: "business shirt",
			Brand: "tommy hilfiger",
			Price: 300,
			Stock: 8,
		},
	}

	for _, p := range resp.Data {
		expectedProduct, ok := products[p.Title]
		if !assert.True(t, ok) {
			t.Logf("Product %q was not expected", p.Title)
			t.FailNow()
		}

		assert.Equal(t, expectedProduct.Title, p.Title)
		assert.Equal(t, expectedProduct.Brand, p.Brand)
		assert.Equal(t, expectedProduct.Price, p.Price)
		assert.Equal(t, expectedProduct.Stock, p.Stock)
	}
}

func TestSearch_FilterByBrand(t *testing.T) {
	ctx := context.Background()
	esClient, cleanup := esClient(t)
	defer cleanup()

	searchSvc := elasticsearch.NewSearchService(esClient)
	searchSvc.InitMapping(ctx)

	err := searchSvc.LoadFile(filepath.Join("testdata", "products.json"))
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	resp, err := searchSvc.Search(ctx, "t-shirt", gfgsearch.Options{
		Page:    1,
		PerPage: 5,
		Filter: map[string]string{
			"brand": "tommy hilfiger",
		},
	})

	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.EqualValues(t, 2, resp.Total)
	assert.Len(t, resp.Data, 2)

	products := map[string]gfgsearch.Product{
		"printed t-shirt": {
			Title: "printed t-shirt",
			Brand: "tommy hilfiger",
			Price: 200,
			Stock: 9,
		},
		"business shirt": {
			Title: "business shirt",
			Brand: "tommy hilfiger",
			Price: 300,
			Stock: 8,
		},
	}

	for _, p := range resp.Data {
		expectedProduct, ok := products[p.Title]
		if !assert.True(t, ok) {
			t.Logf("Product %q was not expected", p.Title)
			t.FailNow()
		}

		assert.Equal(t, expectedProduct.Title, p.Title)
		assert.Equal(t, expectedProduct.Brand, p.Brand)
		assert.Equal(t, expectedProduct.Price, p.Price)
		assert.Equal(t, expectedProduct.Stock, p.Stock)
	}
}

func TestSearch_SortByPriceAsc(t *testing.T) {
	ctx := context.Background()
	esClient, cleanup := esClient(t)
	defer cleanup()

	searchSvc := elasticsearch.NewSearchService(esClient)
	searchSvc.InitMapping(ctx)

	err := searchSvc.LoadFile(filepath.Join("testdata", "products.json"))
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	resp, err := searchSvc.Search(ctx, "", gfgsearch.Options{
		Page:    1,
		PerPage: 5,
		SortBy: map[string]string{
			"price": "asc",
		},
	})

	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.EqualValues(t, 10, resp.Total)
	assert.Len(t, resp.Data, 5)

	for i, p := range resp.Data {
		assert.EqualValues(t, (i+1)*100, p.Price)
	}
}

func TestSearch_SortByStockDesc(t *testing.T) {
	ctx := context.Background()
	esClient, cleanup := esClient(t)
	defer cleanup()

	searchSvc := elasticsearch.NewSearchService(esClient)
	searchSvc.InitMapping(ctx)

	err := searchSvc.LoadFile(filepath.Join("testdata", "products.json"))
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	resp, err := searchSvc.Search(ctx, "", gfgsearch.Options{
		Page:    1,
		PerPage: 5,
		SortBy: map[string]string{
			"stock": "desc",
		},
	})

	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.EqualValues(t, 10, resp.Total)
	assert.Len(t, resp.Data, 5)

	for i, p := range resp.Data {
		assert.EqualValues(t, 10-i, p.Stock)
	}
}
