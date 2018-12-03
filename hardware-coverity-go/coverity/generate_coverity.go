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
	Severity  string `json:"severity"`
	Component string `json:"component"`
}

func SingleTest(raw_result io.ReadCloser, esClient *elastic.Client) (Report, error) {
	// verify the validity of what

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

	sr, err := esClient.Search().Index("hardware").Type("doc").Query(boardQuery).Do(context.Background())

	if err != nil {
		log.Fatal(err)
		return Report{}, err
	}
	hw, err := decodeBoard(sr)

	if err != nil {
		log.Fatal(err)
		return Report{}, err
	}

	fmt.Println(hw)

	var buckets []string

	for i := 0; i < reflect.ValueOf(hw).NumField(); i++ {
		val := reflect.ValueOf(hw).Field(i)
		valType := reflect.TypeOf(hw).Field(i)

		switch t := val.Interface().(type) {
		case bool:
			fmt.Printf("found bool")
			tag := valType.Tag.Get("json")
			if tag != "" {
				buckets = append(buckets, tag)
			}
		case string:
			fmt.Printf("found string")
			tag := val.Interface()
			if tag != "" {
				buckets = append(buckets, tag.(string))
			}
		default:
			fmt.Printf("I don't know about type %T!\n", t)
		}
	}
	fmt.Println(buckets)

	// queries for the related drivers

	return Report{}, nil
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
