package http

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOptionsFrom(t *testing.T) {
	params := make(url.Values)
	params.Set(pageParam, "2")
	params.Set(perPageParam, "10")
	params.Set(filterParam, "title:a,brand:b,price:c,stock:d,invalid=e")
	params.Set(sortParam, "title,-brand,price,-stock,invalid,-wrong")

	opt := getOptionsFrom(params)
	assert.Equal(t, 2, opt.Page)
	assert.Equal(t, 10, opt.PerPage)
	// Filters
	if assert.NotNil(t, opt.Filter["title"]) {
		assert.Equal(t, "a", opt.Filter["title"])
	}
	if assert.NotNil(t, opt.Filter["brand"]) {
		assert.Equal(t, "b", opt.Filter["brand"])
	}
	if assert.NotNil(t, opt.Filter["price"]) {
		assert.Equal(t, "c", opt.Filter["price"])
	}
	if assert.NotNil(t, opt.Filter["stock"]) {
		assert.Equal(t, "d", opt.Filter["stock"])
	}
	assert.NotContains(t, opt.Filter, "invalid")
	// SortBy
	if assert.NotNil(t, opt.SortBy["title"]) {
		assert.Equal(t, "asc", opt.SortBy["title"])
	}
	if assert.NotNil(t, opt.SortBy["brand"]) {
		assert.Equal(t, "desc", opt.SortBy["brand"])
	}
	if assert.NotNil(t, opt.SortBy["price"]) {
		assert.Equal(t, "asc", opt.SortBy["price"])
	}
	if assert.NotNil(t, opt.SortBy["stock"]) {
		assert.Equal(t, "desc", opt.SortBy["stock"])
	}
	assert.NotContains(t, opt.SortBy, "invalid")
	assert.NotContains(t, opt.SortBy, "-wrong")
	assert.NotContains(t, opt.SortBy, "wrong")
}

func TestGetOptionsFrom_DefaultValues(t *testing.T) {
	params := make(url.Values)
	params.Set(pageParam, "invalid")
	params.Set(perPageParam, "invalid")

	opt := getOptionsFrom(params)
	assert.Equal(t, 1, opt.Page)
	assert.Equal(t, DefaultPerPage, opt.PerPage)
}
