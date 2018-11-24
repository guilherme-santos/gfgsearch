package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/guilherme-santos/gfgsearch"
	gfghttp "github.com/guilherme-santos/gfgsearch/http"
	"github.com/guilherme-santos/gfgsearch/mock"

	"github.com/stretchr/testify/assert"
)

func TestSearchHandler(t *testing.T) {
	searcher := mock.NewSearcher(t)
	searcher.SearchFn = func(ctx context.Context, term string, opt gfgsearch.Options) (*gfgsearch.Result, error) {
		assert.Equal(t, "product", term)
		assert.Equal(t, 2, opt.Page)
		assert.Equal(t, 10, opt.PerPage)
		if assert.NotNil(t, opt.Filter["brand"]) {
			assert.Equal(t, "santos", opt.Filter["brand"])
		}
		if assert.NotNil(t, opt.SortBy["price"]) {
			assert.Equal(t, "desc", opt.SortBy["price"])
		}
		if t.Failed() {
			t.FailNow()
		}

		return &gfgsearch.Result{
			Total: 100,
			Data: []gfgsearch.Product{
				{
					Title: "product-a",
					Brand: "brand-a",
					Stock: 12,
					Price: 1234,
				},
			},
		}, nil
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "https://global-fashion-group.com/search?q=product&page=2&per_page=10&filter=brand:santos&sort=-price", nil)

	h := gfghttp.NewSearchHandler(searcher)
	h.ServeHTTP(w, req)

	assert.True(t, searcher.SearchInvoked)
	assert.Equal(t, http.StatusOK, w.Code)

	expectedBody := `{"total":100,"data":[{"title":"product-a","brand":"brand-a","price":1234,"stock":12}]}`
	assert.Equal(t, expectedBody, w.Body.String())
}
