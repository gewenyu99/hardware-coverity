package coverity

import (
	"fmt"
	"github.com/gewenyu99/hardware-coverity/json_requests"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestSingleTest(t *testing.T) {
	example_result_1()
}

func example_result_1() {
	GOPATH := os.Getenv("GOPATH")
	fmt.Println(GOPATH)
	example_result, err := os.Open(strings.Join([]string{GOPATH, "/src/github.com/gewenyu99/hardware-coverity/hardware-coverity-go/sample-data/sample-result-1.json"}, ""))
	if err != nil {
		fmt.Println(err)
	}
	defer example_result.Close()
	example, _ := ioutil.ReadAll(example_result)
	fmt.Println(string(example))
	resp := json_requests.Post("http://localhost:8080/singleTestEntry", string(example))
	if resp == "" {
		panic("Ooooooooooooooooooooopppps")
	}
	fmt.Println(resp)
}