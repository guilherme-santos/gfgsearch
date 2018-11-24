package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/guilherme-santos/gfgsearch"

	"github.com/olivere/elastic"
)

var (
	Index = "gfg_products"
	Type  = "products"
)

type SearchService struct {
	esClient *elastic.Client
}

func NewSearchService(esClient *elastic.Client) *SearchService {
	return &SearchService{
		esClient: esClient,
	}
}

func (s *SearchService) InitMapping(ctx context.Context) error {
	_, err := s.esClient.CreateIndex(Index).
		BodyJson(map[string]interface{}{
			"settings": map[string]interface{}{
				"number_of_shards": 1,
			},
			"mappings": map[string]interface{}{
				Type: map[string]interface{}{
					"properties": map[string]interface{}{
						"title": map[string]interface{}{
							"type": "text",
						},
						"brand": map[string]interface{}{
							"type": "keyword",
						},
						"price": map[string]interface{}{
							"type": "integer",
						},
						"stock": map[string]interface{}{
							"type": "integer",
						},
					},
				},
			},
		}).
		Do(ctx)
	return err
}

func (s *SearchService) LoadFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var products []map[string]interface{}

	err = json.Unmarshal(data, &products)
	if err != nil {
		return err
	}

	bulk := s.esClient.Bulk().Index(Index).Type(Type)

	for _, product := range products {
		bulk.Add(
			elastic.NewBulkIndexRequest().Doc(product),
		)
	}

	ctx := context.Background()

	_, err = bulk.Do(ctx)
	if err != nil {
		return err
	}

	s.esClient.Refresh().Index(Index).Do(ctx)
	return nil
}

func (s *SearchService) Search(ctx context.Context, term string, opt gfgsearch.Options) (*gfgsearch.Result, error) {
	query := elastic.NewBoolQuery()
	if term != "" {
		query.Must(elastic.NewMatchQuery("title", term))
	} else {
		query.Must(elastic.NewMatchAllQuery())
	}

	for field, value := range opt.Filter {
		query.Filter(elastic.NewTermQuery(field, value))
	}

	searchSvc := s.esClient.Search()
	searchSvc.Index(Index)
	searchSvc.Type(Type)
	searchSvc.Timeout(time.Second.String())
	searchSvc.Query(query)
	searchSvc.From((opt.Page - 1) * opt.PerPage)
	searchSvc.Size(opt.PerPage)

	for field, value := range opt.SortBy {
		var ascending bool
		if value == "asc" {
			ascending = true
		}
		searchSvc.Sort(field, ascending)
	}

	esResp, err := searchSvc.Do(ctx)
	if err != nil {
		return nil, err
	}

	var resp gfgsearch.Result

	resp.Total = int32(esResp.TotalHits())
	resp.Data, err = mapProducts(esResp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func mapProducts(esResp *elastic.SearchResult) ([]gfgsearch.Product, error) {
	if esResp == nil || esResp.TotalHits() == 0 {
		return make([]gfgsearch.Product, 0), nil
	}

	products := make([]gfgsearch.Product, 0, esResp.TotalHits())

	for _, hit := range esResp.Hits.Hits {
		var p gfgsearch.Product

		err := json.Unmarshal(*hit.Source, &p)
		if err != nil {
			return nil, fmt.Errorf("cannot unmarshal vendor: %s", err.Error())
		}

		products = append(products, p)
	}

	return products, nil
}
