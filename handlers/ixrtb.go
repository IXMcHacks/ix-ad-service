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
	topBid := processAuction(dspRequest, dspURLs)

	// Return the topBid in JSON format back to the webpage that sent the initial request.
	ReturnJSONResponse(w, topBid)
}

// Public Service Announcement
var fallbackPSAURL = "https://psanycsquad.podbean.com/mf/web/f9rrjh/maxresdefault.jpg"

// Return a Public Service Announcement back to the browser. This is done usually when
// either an error occured during an auction, or if the auction process times out.
func returnPSA(w http.ResponseWriter) {
	ReturnJSONResponse(w, Bid{AdURL: fallbackPSAURL})
}

//////////////////////////////////////////////////////////////////////////////////////////////////
//																								//
//										WORKSHOP COMPONENT										//
//																								//
//		In this workshop you will be completing the processAuction() method. The method			//
//		takes in a DspRequest and an array of URLs to reach DSPs with a http request.			//
//		Your job is to get the top bid from each DSP, compare the bids, and return the  		//
//		highest bid.																			//
//																								//
//////////////////////////////////////////////////////////////////////////////////////////////////

// Replace the following URLs with the URLs provided by the workshop organizers.
var dspURLs = []string{
	"http://10.65.111.204:8080/ixrtb",
	"http://10.65.104.107:8080/ixrtb",
}

// Complete this method
func processAuction(dspRequest DspRequest, dsps []string) Bid {

	// Your job is to populate the topBid object with the highest bid received
	var topBid Bid

	// Some key steps to guide you:

	// 1. Get the bid of that DSP by calling:
	// bid, err := getBidFromDSP(dspRequest, dsps[0]) for example

	// 2. For each bid received, compare it to the highest bid and set the topBid accordingly.

	// 3. HINT: can we find a way to send and receive all requests concurrently using goroutines
	// and channels?

	return topBid
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
