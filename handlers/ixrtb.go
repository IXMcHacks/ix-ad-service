package handlers

import (
	"log"
	"net/http"
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

// RunAuction is the ad-serving handler. It receives an IXRTB GET request, parses the values in the
// url fields, builds an IXRTB POST request to DSPs with those same values a the request body, and
// sends a bid request to DSPs. It then waits for responses, selects the highest bidder and returns
// the ad belonging to the Bid.
func RunAuction(w http.ResponseWriter, r *http.Request) {

	log.Printf("Received Ad-request")

	// Parse the IXRTB GET request and put values into a DspRequest struct
	dspRequest, parseError := parseGETRequest(w, r)
	if parseError != nil {
		ReturnJSONResponse(w, HeartBeatResponse{Okay: false})
	}

	// Instantiate a bidChannel and errChannel through which the bid responses will be passed
	// through. This allows multiple requests to DSPs to be made concurrently, and their responses
	// to be processed as they are received.
	bidChannel := make(chan Bid)
	errChannel := make(chan error)

	// Loop over all the DSPs provided and send a bid request to each of them. A go routine for
	// each DSP is necessary to ensure it is not waiting for the last DSP to respond before
	// contacting another one.
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

	var topBid Bid
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

	ReturnJSONResponse(w, topBid)
}

// getDSPBid is an example of how to build and send a POST request to send to a server.
// It builds the POST request body data using values provided from a DspRequest struct,
// instantiates a http client, and sends and receives the request using it.
func getDSPBid(dspRequest DspRequest, dspURL string) (Bid, error) {

	// Specifiy the variable that will hold the final return result.
	var bid Bid

	// Using the client that was initialized before, invoke the Do method and provide it the
	// POST request that was built. Notice how they are decoupled, so many clients can make
	// the same POST request if they are provided with the same built request. This is also where
	// we receive the response, and assign it to the response variable.
	response, err := sendPOSTRequest(dspRequest, dspURL)
	if err != nil {
		log.Printf("Error at sending request to DSP at: %v", dspURL)
		return bid, err
	}

	return parseJSONResponse(bid, response, dspURL)
}
