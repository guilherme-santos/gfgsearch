package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/guilherme-santos/gfgsearch/elasticsearch"
	gfghttp "github.com/guilherme-santos/gfgsearch/http"
)

func main() {
	srvAddr := os.Getenv("HTTP_SERVER_ADDR")
	esURL := os.Getenv("ELASTICSEARCH_URL")
	log := log.New(os.Stdout, "", log.LstdFlags)

	esClient, err := elasticsearch.NewClient(esURL)
	if err != nil {
		log.Fatalln("Unable to connect to ElasticSearch:", err)
	}

	searchSvc := elasticsearch.NewSearchService(esClient)

	populate := flag.String("populate", "", "use filename provided to populate ElasticSearch")
	flag.Parse()

	if *populate != "" {
		filename := *populate

		err = searchSvc.LoadFile(filename)
		if err != nil {
			fmt.Printf("Unable to populate ElasticSearch: %s\n", err)
			return
		}

		fmt.Printf("File %q loaded successfully\n", filename)
		return
	}

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
