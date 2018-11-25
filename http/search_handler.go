package http

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/guilherme-santos/gfgsearch"
)

const (
	termParam    = "q"
	pageParam    = "page"
	perPageParam = "per_page"
	filterParam  = "filter"
	sortParam    = "sort"
)

// DefaultPerPage will be use when clients do not specify per_page attribute in the
// request.
var DefaultPerPage = 30

type SearchHandler struct {
	searcher gfgsearch.Searcher
}

func NewSearchHandler(s gfgsearch.Searcher) *SearchHandler {
	return &SearchHandler{
		searcher: s,
	}
}

func (h *SearchHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	params := req.URL.Query()
	term := params.Get(termParam)
	opt := getOptionsFrom(params)

	res, err := h.searcher.Search(ctx, term, opt)
	if err != nil {
		Logger.Println("Unable to search for products:", err)
		// TODO add specific errors to response with the right status code.
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(res)
	if err != nil {
		Logger.Println("Unable to marshal response:", err)
		newErrorResponse(w, http.StatusInternalServerError, httpError{
			code:    "invalid_json",
			message: err.Error(),
		})
		return
	}

	newResponse(w, http.StatusOK, body)
}

func getOptionsFrom(params url.Values) gfgsearch.Options {
	opt := gfgsearch.Options{
		Page:    1,
		PerPage: DefaultPerPage,
		Filter:  make(map[string]string),
		SortBy:  make(map[string]string),
	}

	if page, err := strconv.Atoi(params.Get(pageParam)); err == nil {
		opt.Page = page
	}
	if perPage, err := strconv.Atoi(params.Get(perPageParam)); err == nil {
		opt.PerPage = perPage
	}

	filters := strings.Split(params.Get(filterParam), ",")
	for _, f := range filters {
		parts := strings.SplitN(f, ":", 2)
		// if is missing the content ignore the filter
		if len(parts) != 2 || !gfgsearch.IsValidField(parts[0]) {
			continue
		}

		opt.Filter[parts[0]] = parts[1]
	}

	sort := strings.Split(params.Get(sortParam), ",")
	for _, field := range sort {
		var order string

		if strings.HasPrefix(field, "-") {
			field = field[1:]
			order = "desc"
		} else {
			order = "asc"
		}

		if !gfgsearch.IsValidField(field) {
			continue
		}

		opt.SortBy[field] = order
	}

	return opt
}
