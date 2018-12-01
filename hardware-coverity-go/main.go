package main

import "github.com/olivere/elastic"
import "github.com/gewenyu99/hardware-coverity/hardware-coverity-go/sample-data"

func main() {
	esClient, err := elastic.NewClient(elastic.SetURL("http://192.168.2.10:9201"))
	if err != nil {
		// Handle error
	}
	sample_data.LoadSample(esClient)

}
