package main

import (
	"encoding/json"
	"github.com/gewenyu99/hardware-coverity/hardware-coverity-go/coverity"
	"github.com/olivere/elastic"
	"log"
	"net/http"
)
import "github.com/gewenyu99/hardware-coverity/hardware-coverity-go/sample-data"

func main() {
	esClient, err := elastic.NewClient(elastic.SetURL("http://elasticsearch:9200"))
	if err != nil {
		log.Fatal(err)
	}

	sample_data.LoadSample(esClient)

	http.HandleFunc("/singleTestEntry", HandleSingleTestEntry(esClient))

	log.Fatal(http.ListenAndServe(":8080", nil))

}
func HandleSingleTestEntry( esClient *elastic.Client) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		report, err := coverity.SingleTest(r.Body, esClient)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		response, err := json.Marshal(report)
		http.Error(w, string(response), 200)
	})
}
