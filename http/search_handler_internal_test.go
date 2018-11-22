package http

import (
	"net/url"
	"testing"
)

func TestGetOptionsFrom(t *testing.T) {
	params := make(url.Values)
	params.Set(pageParam, "2")
	params.Set(perPageParam, "10")
	params.Set(filterParam, "title:a,brand:b,price:c,stock:d,invalid=e")
	params.Set(sortParam, "title,-brand,price,-stock,invalid,-wrong")

	opt := getOptionsFrom(params)
	if opt.Page != 2 {
		t.Fatalf("Expected page to be %d but got %d", 2, opt.Page)
	}
	if opt.PerPage != 10 {
		t.Fatalf("Expected per_page to be %d but got %d", 10, opt.PerPage)
	}
	// Filters
	if title, ok := opt.Filter["title"]; !ok || title != "a" {
		t.Fatalf("Expected filter title to be %s but got %s", "a", title)
	}
	if brand, ok := opt.Filter["brand"]; !ok || brand != "b" {
		t.Fatalf("Expected filter brand to be %s but got %s", "b", brand)
	}
	if price, ok := opt.Filter["price"]; !ok || price != "c" {
		t.Fatalf("Expected filter price to be %s but got %s", "c", price)
	}
	if stock, ok := opt.Filter["stock"]; !ok || stock != "d" {
		t.Fatalf("Expected filter stock to be %s but got %s", "d", stock)
	}
	if _, ok := opt.Filter["invalid"]; ok {
		t.Fatal("Field invalid is not searchable and shouldn't be in the filter")
	}
	// SortBy
	if order, ok := opt.SortBy["title"]; !ok || order != "asc" {
		t.Fatalf("Expected title order to be %s but got %s", "asc", order)
	}
	if order, ok := opt.SortBy["brand"]; !ok || order != "desc" {
		t.Fatalf("Expected brand order to be %s but got %s", "asc", order)
	}
	if order, ok := opt.SortBy["price"]; !ok || order != "asc" {
		t.Fatalf("Expected price order to be %s but got %s", "asc", order)
	}
	if order, ok := opt.SortBy["stock"]; !ok || order != "desc" {
		t.Fatalf("Expected stock order to be %s but got %s", "asc", order)
	}
	if _, ok := opt.SortBy["invalid"]; ok {
		t.Fatal("Field invalid is not searchable and shouldn't be in the sort")
	}
	if _, ok := opt.SortBy["-wrong"]; ok {
		t.Fatal("Field -wrong is not searchable and shouldn't be in the sort")
	}
	if _, ok := opt.SortBy["wrong"]; ok {
		t.Fatal("Field wrong is not searchable and shouldn't be in the sort")
	}
}

func TestGetOptionsFrom_DefaultValues(t *testing.T) {
	params := make(url.Values)
	params.Set(pageParam, "invalid")
	params.Set(perPageParam, "invalid")

	opt := getOptionsFrom(params)
	if opt.Page != 1 {
		t.Fatalf("Expected page to be %d but got %d", 1, opt.Page)
	}
	if opt.PerPage != DefaultPerPage {
		t.Fatalf("Expected per_page to be %d but got %d", DefaultPerPage, opt.PerPage)
	}
}
