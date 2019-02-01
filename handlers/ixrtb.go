package handlers

import (
	"log"
	"net/http"
)

// Bid is the stuct which will hold the information regarding an advertiser's bid for their
// ad to be shown on the available ad-placement.
type Bid struct {
	Advertiser    string `json:"Advertiser"`
	BidPrice      int    `json:"BidPrice"`
	AdURL         string `json:"AdURL"`
	AdDescription string `json:"AdDescription"`
}

// DspRequest is the IXRTB standard compliant struct that is used to exchange information
// about a particular ad-placement. It contains the fields "Size", "Likes", "Dislikes", "Age",
// and "Code", which advertisers use to determine how much they should bit for their ad to fill
// the placement.
type DspRequest struct {
	S    string `schema:"s"`
	L    string `schema:"l"`
	D    string `schema:"d"`
	A    int    `schema:"a"`
	Code string `schema:"code"`
}

// RunAuction is the ad-serving handler. It receives an IXRTB GET request, parses the values in the
// url fields, builds an IXRTB POST request to DSPs with those same values a the request body, and
// sends a bid request to DSPs. It then waits for responses, selects the highest bidder and returns
// the ad belonging to the Bid.
func RunAuction(w http.ResponseWriter, r *http.Request) {

	log.Printf("Received Ad-request!")

	// Parse the IXRTB GET request and put values into a DspRequest struct.
	dspRequest, parseError := parseGETRequest(w, r)
	if parseError != nil {
		returnPSA(w)
	}

	// We call processAuction to get the topBid whose ad we will return to the browser.
	topBid := processAuction(dspRequest)

	// Return the topBid in JSON format back to the webpage that sent the initial request.
	ReturnJSONResponse(w, topBid)
}

// getBidFromDSP sends a bid request to a specific DSP, parses its response, an returns
// the parsed response as fully useable Bid object that can be processed.
func getBidFromDSP(dspRequest DspRequest, dspURL string) (Bid, error) {

	// Specifiy the variable that will hold the final return result.
	var bid Bid

	// Call the helper sendPOSTRequest method to send the post request.
	response, err := sendPOSTRequest(dspRequest, dspURL)
	if err != nil {
		log.Printf("Error at sending request to DSP at: %v", dspURL)
		return bid, err
	}

	// Call the helper parseJSONResponse method to parse the response received.
	return parseJSONResponse(bid, response, dspURL)
}

// Public Service Announcement
type fallbackPSA struct {
	AdURL string `json:"AdURL"`
}

var fallbackPSAURL = "https://psanycsquad.podbean.com/mf/web/f9rrjh/maxresdefault.jpg"

// Return a Public Service Announcement back to the browser. This is done usually when
// either an error occured during an auction, or if the auction process times out.
func returnPSA(w http.ResponseWriter) {
	ReturnJSONResponse(w, fallbackPSA{AdURL: fallbackPSAURL})
}

//////////////////////////////////////////////////////////////////////////////////////////////////

//										WORKSHOP COMPONENT										//

//////////////////////////////////////////////////////////////////////////////////////////////////

// dspURLs is the list of urls to send requests to DSPs
var dspURLs = []string{
	"http://10.65.111.204:8080/ixrtb",
	"http://10.65.104.107:8080/ixrtb",
}

// processAuction is what processes the sending of bid-requests to DSPs, receiving, parsing and
// validating their responses, and returning the highest bid.
func processAuction(dspRequest DspRequest) Bid {

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
			bid, err := getBidFromDSP(dspRequest, dspURL)
			if err != nil {
				errChannel <- err
			} else {
				bidChannel <- bid
			}
		}(dspURL)
	}

	// The ad with the topBid that will be returned to the website is prepared.
	var topBid Bid

	// Loop continuously and listen for responses from DSPs. The select block waits for
	// the channels to receive data and runs the code according to which channel it received
	// data from. As long as not all responses are received, the select block will keep waiting
	// until one of the cases are triggered.
	// When a bid is received, we check it against the current highest bid and determine the new
	// highest bid. This repeats until all expected responses are received.
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

	return topBid
}
