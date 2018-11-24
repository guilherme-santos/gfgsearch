package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	gfghttp "github.com/guilherme-santos/gfgsearch/http"

	"github.com/stretchr/testify/assert"
)

func TestBasicAuthMiddleware(t *testing.T) {
	var handlerCalled bool

	var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusCreated)
	})
	h = gfghttp.BasicAuthMiddleware(h, "search", "gfc")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "https://global-fashion-group.com/search", nil)
	req.SetBasicAuth("search", "gfc")

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.True(t, handlerCalled)
}

func TestBasicAuthMiddleware_NoAuthentication(t *testing.T) {
	var handlerCalled bool

	var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusCreated)
	})
	h = gfghttp.BasicAuthMiddleware(h, "search", "gfc")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "https://global-fashion-group.com/search", nil)

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.False(t, handlerCalled)
}

func TestBasicAuthMiddleware_WrongPassword(t *testing.T) {
	var handlerCalled bool

	var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusCreated)
	})
	h = gfghttp.BasicAuthMiddleware(h, "search", "gfc")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "https://global-fashion-group.com/search", nil)
	req.SetBasicAuth("search", "wrong")

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.False(t, handlerCalled)
}
