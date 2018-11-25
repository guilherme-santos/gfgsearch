package mock

import (
	"context"
	"testing"

	"github.com/guilherme-santos/gfgsearch"
)

// Searcher is a mock implementation of gfgsearch.Searcher.
type Searcher struct {
	t *testing.T

	// SearchInvoked will be set as true if method Search() is called.
	SearchInvoked bool
	// SearchFn should be set with your own implementation before you
	// use this mock.
	SearchFn func(ctx context.Context, term string, opt gfgsearch.Options) (*gfgsearch.Result, error)
}

func NewSearcher(t *testing.T) *Searcher {
	return &Searcher{
		t: t,
	}
}

func (s *Searcher) Search(ctx context.Context, term string, opt gfgsearch.Options) (*gfgsearch.Result, error) {
	if s.SearchFn == nil {
		s.t.Fatal("You need to set SearchFn to use this mock")
	}

	s.SearchInvoked = true
	return s.SearchFn(ctx, term, opt)
}
