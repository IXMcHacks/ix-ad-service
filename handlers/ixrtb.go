package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/schema"
)

type Bid struct {
	Advertiser    string `json:"Advertiser"`
	BidPrice      int    `json:"BidPrice"`
	AdURL         string `json:"AdURL"`
	AdDescription string `json:"AdDescription"`
}

type DspRequest struct {
	S    string `schema:"s"`
	L    string `schema:"l"`
	D    string `schema:"d"`
	A    int    `schema:"a"`
	Code string `schema:"code"`
}

type HeartBeatResponse struct {
	Okay bool `json:"okay"`
}

var dspURLs = []string{
	"http://127.0.0.1:9000/ixrtb",
	"http://127.0.0.1:9002/ixrtb",
}

// ServeAd is the ad-serving handler. It receives an IXRTB GET request, parses the values in the
// url fields, builds an IXRTB POST request to DSPs with those same values a the request body, and
// sends a bid request to DSPs. It then waits for responses, selects the highest bidder and returns
// the ad belonging to the Bid.
func RunAuction(w http.ResponseWriter, r *http.Request) {

	log.Printf("Received Ad-request")

	// Parse the IXRTB GET request and put values into a DspRequest struct
	dspRequest, parseError := parseGETRequest(w, r)
	if parseError != nil {
		HandleSuccess(&w, HeartBeatResponse{Okay: false})
	}

	var topBid Bid

	bidChannel := make(chan Bid)
	errChannel := make(chan error)

	for _, dspURL := range dspURLs {
		go func(dspURL string) {
			bid, err := getDSPBid(dspRequest, dspURL)
			if err != nil {
				errChannel <- err
			} else {
				bidChannel <- bid
			}
		}(dspURL)
	}

	for responsesReceived := 0; responsesReceived < len(dspURLs); {
		select {
		case gotBid := <-bidChannel:
			if gotBid.BidPrice > topBid.BidPrice {
				topBid = gotBid
			}
			responsesReceived++
		case gotError := <-errChannel:
			log.Printf("Error getting request from DSP: %v", gotError)
			responsesReceived++
		}
	}

	HandleSuccess(&w, topBid)
}

func parseGETRequest(w http.ResponseWriter, r *http.Request) (DspRequest, error) {

	// Initialize DspRequest struct we want to populate
	var dspRequest DspRequest

	var decoder = schema.NewDecoder()
	values := r.URL.Query()
	err := decoder.Decode(&dspRequest, values)
	if err != nil {
		log.Printf("Error in GET parameters : %v", err)
		return dspRequest, err
	}

	return dspRequest, nil
}

func getDSPBid(dspRequest DspRequest, dspURL string) (Bid, error) {

	data := url.Values{}
	data.Set("s", dspRequest.S)
	data.Set("l", dspRequest.L)
	data.Set("d", dspRequest.D)
	data.Set("a", strconv.Itoa(dspRequest.A))

	client := &http.Client{}

	dspRequestBody, _ := http.NewRequest("POST", dspURL, strings.NewReader(data.Encode()))

	dspRequestBody.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	dspRequestBody.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	log.Printf("Sending request to dsp:%v", dspURL)
	response, err := client.Do(dspRequestBody)

	var bid Bid

	if err != nil {
		log.Printf("Error at sending request to DSP:%v", dspURL)
		return bid, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error at response body from DSP:%v", dspURL)
		return bid, err
	}

	err = json.Unmarshal(body, &bid)
	if err != nil {
		log.Printf("Unable to unmarshal the json body of DSP:%v", dspURL)
		return bid, err
	}

	return bid, nil
}
