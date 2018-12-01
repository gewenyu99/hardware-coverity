package sample_data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	"io/ioutil"
	"os"
	"strings"
)

type TestResult struct {
	Pass bool `json:"pass"`
	config TestConfig `json:"config"`
}

type TestConfig struct {
	Target string `json:"target"`
	EnabledIO []string `json:"enabled_io"`
}

type Hardware struct {
	BoardName            string `json:"board-name"`
	BoardType            string `json:"board-type"`
	Architecture         string `json:"architecture"`
	Graphics             string `json:"graphics"`
	Usb3                 bool   `json:"usb3"`
	Usb2                 bool   `json:"usb2"`
	Gigabit              bool   `json:"gigabit"`
	IntegratedNetworking bool   `json:"integrated_networking"`
	HDMI                 bool   `json:"hdmi"`
	SATA                 bool   `json:"sata"`
	Mdot2                bool   `json:"m.2"`
	Optane               bool   `json:"optane"`
	ECC                  bool   `json:"ecc"`
	VirtIO               bool   `json:"virtio"`
	Raid                 bool   `json:"raid"`
	Audio                bool   `json:"audio"`
}

type HardwareList struct {
	Hardware []Hardware `json:"harwareList"`
}

type Driver struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	DeviceType string   `json:"device-type"`
	Tags       []string `json:"tags"`
}

type FitterPacks struct {
	FitterPack []Driver `json:"fitter-pack"`
}

func LoadSample(esClient *elastic.Client) {
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

func loadSampleDrivers() (error, FitterPacks) {
	GOPATH := os.Getenv("GOPATH")
	fmt.Println(GOPATH)
	driversFile, err := os.Open(strings.Join([]string{GOPATH,"/src/github.com/gewenyu99/hardware-coverity/hardware-coverity-go/sample-data/sample-drivers.json"} ,""))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened sample-drivers.json")
	defer driversFile.Close()
	drivesByte, _ := ioutil.ReadAll(driversFile)
	var drivers FitterPacks
	err = json.Unmarshal(drivesByte, &drivers)
	if err != nil {
		fmt.Println(err)
	}
	return err, drivers
}

func loadSampleHardware() (error, HardwareList) {
	GOPATH := os.Getenv("GOPATH")
	fmt.Println(GOPATH)
	hardwareFile, err := os.Open(strings.Join([]string{GOPATH,"/src/github.com/gewenyu99/hardware-coverity/hardware-coverity-go/sample-data/sample-hardware.json"} ,""))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened sample-hardware.json")
	defer hardwareFile.Close()
	hardwareBytes, _ := ioutil.ReadAll(hardwareFile)
	var hardware HardwareList
	err = json.Unmarshal(hardwareBytes, &hardware)
	if err != nil {
		fmt.Println(err)
	}
	return err, hardware
}