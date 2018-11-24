package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/guilherme-santos/gfgsearch/elasticsearch"
	gfghttp "github.com/guilherme-santos/gfgsearch/http"

	"github.com/olivere/elastic"
)

func main() {
	srvAddr := os.Getenv("HTTP_SERVER_ADDR")
	esURL := os.Getenv("ELASTICSEARCH_URL")
	log := log.New(os.Stdout, "", log.LstdFlags)

	httpClient := &http.Client{
		Timeout: 2 * time.Second,
	}

	esClient, err := elastic.NewClient(
		elastic.SetURL(esURL),
		elastic.SetSniff(false),
		elastic.SetHttpClient(httpClient),
	)
	if err != nil {
		log.Fatalln("Unable to connect to ElasticSearch:", err)
	}

	searchSvc := elasticsearch.NewSearchService(esClient)
	searchHandler := gfghttp.NewSearchHandler(searchSvc)

	serverMux := http.NewServeMux()
	serverMux.Handle("/v1/search/products",
		gfghttp.LogMiddleware(
			gfghttp.BasicAuthMiddleware(searchHandler, "gfg", "search"),
		),
	)

	srv := &http.Server{
		Handler: serverMux,
		Addr:    srvAddr,
	}
	log.Println("Running http server on", srv.Addr)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalln("Unable to run http server:", err)
	}
}
