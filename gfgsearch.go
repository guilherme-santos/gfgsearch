package gfgsearch

import "context"

// Product contains all the information about a specific product.
type Product struct {
	Title string `json:"title"`
	Brand string `json:"brand"`
	Price int32  `json:"price"`
	Stock int32  `json:"stock"`
}

// Options contains some properties that can be available or not in
// the request and it'll be used during the search phase.
type Options struct {
	Page    int
	PerPage int
	// Filter is a map that contains the field name and the text to filter.
	Filter map[string]string
	// SortBy is a map that contains the field name and the order.
	SortBy map[string]string
}

type Result struct {
	Total int32     `json:"total"`
	Data  []Product `json:"data"`
}

type Searcher interface {
	Search(ctx context.Context, term string, opt Options) (*Result, error)
}

// IsFieldSearchable checks if the field can be filtered (if it's a valid field).
func IsFieldSearchable(field string) bool {
	switch field {
	case "title", "brand", "price", "stock":
		return true
	}
	return false
}
