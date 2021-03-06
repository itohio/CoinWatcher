// Package util supports specific parsing
package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hexoul/go-coinmarketcap/types"
)

// ToInt helper for parsing strings to int
func ToInt(rawInt string) int {
	parsed, _ := strconv.Atoi(strings.Replace(strings.Replace(rawInt, "$", "", -1), ",", "", -1))
	return parsed
}

// ToFloat helper for parsing strings to float
func ToFloat(rawFloat string) float64 {
	parsed, _ := strconv.ParseFloat(strings.Replace(strings.Replace(strings.Replace(rawFloat, "$", "", -1), ",", "", -1), "%", "", -1), 64)
	return parsed
}

// ParseOptions returns parsed param list
func ParseOptions(options *types.Options) (params []string) {
	if options == nil {
		return
	}
	bOption, err := json.Marshal(options)
	if err != nil {
		return
	}
	mOption := new(map[string]interface{})
	if err := json.Unmarshal(bOption, &mOption); err != nil {
		return
	}
	for k, v := range *mOption {
		params = append(params, fmt.Sprintf("%s=%v", k, v))
	}
	return
}

// DoReq HTTP client
func DoReq(req *http.Request) (body []byte, err error) {
	requestTimeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: requestTimeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	}
	return
}

// ISODate returns ISO date format like "2018-11-23"
func ISODate(t time.Time) string {
	return strings.Split(t.Format(time.RFC3339), "T")[0]
}
