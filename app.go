package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/prometheus/common/model"
)

type ErrorType string

type ApiResponse struct {
	Status    string          `json:"status"`
	Data      json.RawMessage `json:"data"`
	ErrorType ErrorType       `json:"errorType,omitempty"`
	Error     string          `json:"error,omitempty"`
}

type QueryResult struct {
	Type   model.ValueType `json:"resultType"`
	Result json.RawMessage `json:"result"`

	// The decoded value.
	//	v model.Value
}

type LabelResult []*string

//Define a new structure that represents out API response (response status and body)
type HTTPResponse struct {
	status string
	//matrix model.Matrix
	result []byte
}

func (r ApiResponse) Successful() bool {
	return string(r.Status) == "success"
}

func main() {
	h := http.HandlerFunc(proxyHandler)
	log.Fatal(http.ListenAndServe(":6789", h))
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received:", r.URL)
	urlParts := strings.Split(r.RequestURI, `/`)
	api := strings.Split(urlParts[3], `?`)[0]

	var ch chan HTTPResponse = make(chan HTTPResponse)
	servers := os.Args[1:]
	for _, server := range servers {
		go DoHTTPGet(fmt.Sprintf("%s%s", server, r.URL), ch)
	}

	remote := make([]ApiResponse, 0)
	for range servers {
		remote = append(remote, ParseResponse((<-ch).result))
	}
	merged := Merge(api, remote)
	asJson, err := json.Marshal(merged)
	if err != nil {
		log.Println("error marshalling back", err)
	}
	fmt.Fprintf(w, "%s", asJson)
}

func DoHTTPGet(url string, ch chan<- HTTPResponse) {
	//Execute the HTTP get
	log.Println("Getting", url)
	httpResponse, _ := http.Get(url)
	httpBody, _ := ioutil.ReadAll(httpResponse.Body)
	//Send an HTTPResponse back to the channel
	ch <- HTTPResponse{httpResponse.Status, httpBody}
}

func ParseResponse(body []byte) ApiResponse {
	var ar ApiResponse
	err := json.Unmarshal(body, &ar)

	if err != nil {
		log.Println("Error unmarshalling JSON api response", err, body)
		return ApiResponse{Status: "error"}
	}

	return ar
}

func Merge(api string, responses []ApiResponse) ApiResponse {
	var qr QueryResult
	var ar ApiResponse

	if len(responses) > 1 {
		ar = responses[0]
	} else {
		log.Println("No responses received")
		return ApiResponse{}
	}

	switch api {
	case `label`:
		return MergeArrays(responses)
	case `series`:
		return MergeSeries(responses)
	case `query_range`:
		err := json.Unmarshal(ar.Data, &qr)
		if err != nil {
			log.Println("Error unmarshalling JSON query result", err, string(ar.Data))
			log.Println("Full response:", string(ar.Data))
			return ApiResponse{}
		}

		switch qr.Type {
		case model.ValMatrix:
			return MergeMatrices(responses)
		default:
			log.Println("Did not recognize the response type of", qr.Type)

		}
	}

	log.Println("Full response:", string(ar.Data))
	return ApiResponse{}
}

func MergeSeries(responses []ApiResponse) ApiResponse {
	merged := make([]map[string]string, 0)

	for _, ar := range responses {
		var result []map[string]string
		err := json.Unmarshal(ar.Data, &result)
		if err != nil {
			log.Println("Unmarshal problem got", err)
		}
		for _, r := range result {
			merged = append(merged, r)
		}
	}
	log.Println("result", merged)

	m, err := json.Marshal(merged)
	if err != nil {
		log.Println("error marshalling series back", err)
	}

	return ApiResponse{Status: "success", Data: m}
}

func MergeArrays(responses []ApiResponse) ApiResponse {
	set := make(map[string]struct{})

	for _, ar := range responses {
		var labels LabelResult
		err := json.Unmarshal(ar.Data, &labels)
		if err != nil {
			log.Println("Error unmarshalling labels", err)
		}
		for _, label := range labels {
			set[*label] = struct{}{}
		}
	}
	keys := make([]string, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	m, err := json.Marshal(keys)
	if err != nil {
		log.Println("error marshalling labels back", err)
	}

	return ApiResponse{Status: "success", Data: m}

}

func MergeMatrices(responses []ApiResponse) ApiResponse {
	samples := make([]*model.SampleStream, 0)
	for _, r := range responses {
		matrix := ExtractMatrix(r)
		for _, s := range matrix {
			samples = append(samples, s)
		}
	}

	mj, _ := json.Marshal(samples)

	qr := QueryResult{model.ValMatrix, mj}
	qrj, _ := json.Marshal(qr)

	r := ApiResponse{Status: "success", Data: qrj}
	return r

}

func ExtractMatrix(ar ApiResponse) model.Matrix {
	var qr QueryResult
	err := json.Unmarshal(ar.Data, &qr)

	if err != nil {
		log.Println("Error unmarshalling JSON query result", err, string(ar.Data))
	}

	var m model.Matrix
	err = json.Unmarshal(qr.Result, &m)
	if err != nil {
		log.Println("Error unmarshalling a matrix", err, string(qr.Result))
	}
	return m
}
