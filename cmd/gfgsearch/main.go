package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/guilherme-santos/gfgsearch/elasticsearch"
	"github.com/guilherme-santos/gfgsearch/http"
)

func main() {
	srvAddr := os.Getenv("HTTP_SERVER_ADDR")
	srvAuthUser := os.Getenv("HTTP_SERVER_AUTH_USER")
	srvAuthPasswd := os.Getenv("HTTP_SERVER_AUTH_PASSWD")
	esURL := os.Getenv("ELASTICSEARCH_URL")
	log := log.New(os.Stdout, "", log.LstdFlags)

	esClient, err := elasticsearch.NewClient(esURL)
	if err != nil {
		log.Fatalln("Unable to connect to ElasticSearch:", err)
	}

	searchSvc := elasticsearch.NewSearchService(esClient)
	// Be sure always when the app run we have mapping created.
	searchSvc.InitMapping(context.Background())

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

	searchHandler := http.NewSearchHandler(searchSvc)

	srv := http.NewServer()
	srv.Use(http.LogMiddleware())
	srv.Use(http.BasicAuthMiddleware(srvAuthUser, srvAuthPasswd))
	srv.Handle("/v1/search/products", searchHandler)

	log.Println("Running http server on", srvAddr)

	err = srv.Listen(srvAddr)
	if err != nil {
		log.Fatalln("Unable to run http server:", err)
	}
}
