package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/guilherme-santos/gfgsearch"
	gfghttp "github.com/guilherme-santos/gfgsearch/http"
	"github.com/guilherme-santos/gfgsearch/mock"
)

func TestSearchHandler(t *testing.T) {
	searcher := mock.NewSearcher(t)
	searcher.SearchFn = func(ctx context.Context, term string, opt gfgsearch.Options) (*gfgsearch.Result, error) {
		if !strings.EqualFold("product", term) {
			t.Fatalf("Expected term to be %s but got %s", "product", term)
		}
		if opt.Page != 2 {
			t.Fatalf("Expected page to be %d but got %d", 2, opt.Page)
		}
		if opt.PerPage != 10 {
			t.Fatalf("Expected per_page to be %d but got %d", 10, opt.PerPage)
		}
		if brand, ok := opt.Filter["brand"]; !ok || brand != "santos" {
			t.Fatalf("Expected filter brand to be %s but got %s", "a", brand)
		}
		if order, ok := opt.SortBy["price"]; !ok || order != "desc" {
			t.Fatalf("Expected price order to be %s but got %s", "asc", order)
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

	if !searcher.SearchInvoked {
		t.Fatal("Expected to call searcher.Search")
	}
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code to be %d but got %d", http.StatusOK, w.Code)
	}
	expectedBody := `{"total":100,"data":[{"title":"product-a","brand":"brand-a","price":1234,"stock":12}]}`
	if !strings.EqualFold(expectedBody, w.Body.String()) {
		t.Fatalf("Expected body to be %q but got %q", expectedBody, w.Body.String())
	}
}
