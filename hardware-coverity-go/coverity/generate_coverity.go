package coverity

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"reflect"
)

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

type Report struct {
	Count int    `json:"count"`
	Risks []Risk `json:"risks"`
}

type TestConfig struct {
	Target    string   `json:"target"`
	EnabledIO []string `json:"enabled_io"`
}

type TestResult struct {
	Pass   bool       `json:"pass"`
	Config TestConfig `json:"config"`
}

type Risk struct {
	Type      string `json:"type"`
	Severity  float64 `json:"severity"`
	Component string `json:"component"`
}

func SingleTest(raw_result io.ReadCloser, esClient *elastic.Client) (Report, error) {
	// verify the validity of what

	fmt.Println("Not on dep!")

	var result TestResult
	byteResult, err := ioutil.ReadAll(raw_result)
	if err != nil {
		log.Fatal(err)
		return Report{}, err
	}
	err = json.Unmarshal(byteResult, &result)

	if err != nil {
		log.Fatal(err)
		return Report{}, err
	}

	// queries for the board in question

	boardQuery := elastic.NewTermQuery("board-name.keyword", result.Config.Target)

	searchResult, err := esClient.Search().Index("hardware").Type("doc").Query(boardQuery).Do(context.Background())

	if err != nil {
		log.Fatal(err)
		return Report{}, err
	}
	hw, err := decodeBoard(searchResult)

	if err != nil {
		log.Fatal(err)
		return Report{}, err
	}

	var scoring_tags []string

	scoring_tags = append(scoring_tags, result.Config.EnabledIO...)

	for i := 0; i < reflect.ValueOf(hw).NumField(); i++ {
		val := reflect.ValueOf(hw).Field(i)
		valType := reflect.TypeOf(hw).Field(i).Name
		switch t := val.Interface().(type) {
		case string:
			if valType != "name" && valType != "BoardName" {
				tag := val.Interface()
				if tag != "" {
					scoring_tags = append(scoring_tags, tag.(string))
				}
			}

		default:
			log.Println("the following was ignored: ", t, val)
		}
	}

	risks, err := findRisks(esClient, scoring_tags)
	return Report{len(risks), risks}, nil
}

func findRisks(esClient *elastic.Client, scoring_tags []string) ([]Risk, error){
	//riskQuery := elastic.NewBoolQuery()

	var queries []elastic.Query
	for _, tag := range scoring_tags {
		fmt.Println(tag)
		queries = append(queries, elastic.NewTermQuery("tags", tag))
	}
	searchResult, err := esClient.Search().Index("driver").Query(elastic.NewBoolQuery().Should(queries...)).Do(context.Background())
	if err != nil {
		return nil, err
	}
	d, err := decodeFitterPack(searchResult)
	fmt.Println(d)
	if err != nil {
		return nil, err
	}

	fmt.Println("Found risks with these tags: ", scoring_tags)
	return d, nil
}

func decodeBoard(res *elastic.SearchResult) (Hardware, error) {
	if res == nil || res.TotalHits() == 0 {
		return Hardware{}, errors.New("The test target doesn't exist")
	}

	var hw Hardware
	err := json.Unmarshal(*res.Hits.Hits[0].Source, &hw)
	if err != nil {
		log.Fatal(err)
		return Hardware{}, err
	}
	return hw, nil
}

func decodeFitterPack(res *elastic.SearchResult) ([]Risk, error) {
	if res == nil || res.TotalHits() == 0 {
		return []Risk{}, errors.New("The test fitter pack doesn't exist")
	}

	var d []Risk

	for i := int64(0); i < res.Hits.TotalHits; i ++{
		temp := Driver{}
		err := json.Unmarshal(*res.Hits.Hits[i].Source, &temp)

		if err != nil {
			log.Fatal(err)
			return []Risk{}, err
		}
		d = append(d, Risk{Type: temp.Type, Severity: *res.Hits.Hits[i].Score, Component:temp.Name})
	}
	return d, nil
}

