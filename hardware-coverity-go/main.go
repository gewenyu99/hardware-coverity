package main

import (
	"fmt"
	"github.com/olivere/elastic"
	"html"
	"log"
	"net/http"
)
import "github.com/gewenyu99/hardware-coverity/hardware-coverity-go/sample-data"

func main() {
	esClient, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	if err != nil {
		// Handle error
	}
	sample_data.LoadSample(esClient)

	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
