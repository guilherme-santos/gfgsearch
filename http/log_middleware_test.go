package http_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	gfghttp "github.com/guilherme-santos/gfgsearch/http"

	"github.com/stretchr/testify/assert"
)

func TestLogMiddleware(t *testing.T) {
	var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})
	h = gfghttp.LogMiddleware(h)

	var buf bytes.Buffer
	gfghttp.Logger = log.New(&buf, "", 0)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "https://global-fashion-group.com/search", nil)
	req.Header.Set("User-Agent", "GDF-HttpClient/1.0")
	req.RemoteAddr = "127.0.0.1:1234"
	req.RequestURI = req.URL.RequestURI()

	h.ServeHTTP(w, req)

	logRegex := regexp.MustCompile(`127\.0\.0\.1 \[[0-9]+(.[0-9]+)?[nµm]?s\] "GET /search HTTP/1.1" 201 Created "GDF-HttpClient/1.0"`)
	assert.Regexp(t, logRegex, buf.String())
}
