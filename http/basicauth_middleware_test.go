package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	gfghttp "github.com/guilherme-santos/gfgsearch/http"
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

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status code %d but got %d", http.StatusCreated, w.Code)
	}
	if !handlerCalled {
		t.Fatal("Dummy handler was expected to be called")
	}
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

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status code %d but got %d", http.StatusUnauthorized, w.Code)
	}
	if handlerCalled {
		t.Fatal("Dummy handler was not expected to be called")
	}
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

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status code %d but got %d", http.StatusUnauthorized, w.Code)
	}
	if handlerCalled {
		t.Fatal("Dummy handler was not expected to be called")
	}
}
