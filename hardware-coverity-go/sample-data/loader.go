package sample_data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gewenyu99/hardware-coverity/hardware-coverity-go/coverity"
	"github.com/olivere/elastic"
	"io/ioutil"
	"os"
	"strings"
)

func LoadSample(esClient *elastic.Client) {
	esClient.DeleteIndex("driver").Do(context.Background())
	esClient.DeleteIndex("hardware").Do(context.Background())
	err, drivers := loadSampleDrivers()
	err, hardwares := loadSampleHardware()

	bulkService := esClient.Bulk()

	for _, driver := range drivers.FitterPack {
		fmt.Println(driver)
		bulkService = bulkService.Add(elastic.NewBulkIndexRequest().Index("driver").Type("doc").Doc(driver))
	}
	for _, hardware := range hardwares.Hardware {
		fmt.Println(hardware)
		bulkService = bulkService.Add(elastic.NewBulkIndexRequest().Index("hardware").Type("doc").Doc(hardware))
	}

	_, err = bulkService.Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}

func loadSampleDrivers() (error, coverity.FitterPacks) {
	GOPATH := os.Getenv("GOPATH")
	fmt.Println(GOPATH)
	driversFile, err := os.Open(strings.Join([]string{GOPATH, "/src/github.com/gewenyu99/hardware-coverity/hardware-coverity-go/sample-data/sample-drivers.json"}, ""))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened sample-drivers.json")
	defer driversFile.Close()
	drivesByte, _ := ioutil.ReadAll(driversFile)
	var drivers coverity.FitterPacks
	err = json.Unmarshal(drivesByte, &drivers)
	if err != nil {
		fmt.Println(err)
	}
	return err, drivers
}

func loadSampleHardware() (error, coverity.HardwareList) {
	GOPATH := os.Getenv("GOPATH")
	fmt.Println(GOPATH)
	hardwareFile, err := os.Open(strings.Join([]string{GOPATH, "/src/github.com/gewenyu99/hardware-coverity/hardware-coverity-go/sample-data/sample-hardware.json"}, ""))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened sample-hardware.json")
	defer hardwareFile.Close()
	hardwareBytes, _ := ioutil.ReadAll(hardwareFile)
	var hardware coverity.HardwareList
	err = json.Unmarshal(hardwareBytes, &hardware)
	if err != nil {
		fmt.Println(err)
	}
	return err, hardware
}
