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

// HTTP Response Header Codes
const (
	httpOK                  = 200
	httpInternalServerError = 500
)

// parseGETRequest is an example of how to read an incoming GET request and extract the
// data in its url into a dedicated struct that can be used by other Go methods in the
// project.
func parseGETRequest(w http.ResponseWriter, r *http.Request) (DspRequest, error) {

	// Initialize DspRequest struct we want to populate with data.
	var dspRequest DspRequest

	// Package gorilla/schema converts structs to and from form values.
	// Here a decoder is instantiated to be able to read the form values
	// of the incoming get request.
	var decoder = schema.NewDecoder()

	// The values in the request URL are extracted using Query() and are converted
	// to form values.
	values := r.URL.Query()

	// The decoder can now insert the form values into the DspRequest struct.
	err := decoder.Decode(&dspRequest, values)
	if err != nil {
		log.Printf("Error in GET parameters : %v", err)
		return dspRequest, err
	}

	return dspRequest, nil
}

// sendPOSTRequest is an example of how to send a POST request to a server and receive
// back a response.
func sendPOSTRequest(dspRequest DspRequest, dspURL string) (*http.Response, error) {

	// 1. Build the POST Request

	// First we must get the data from the DspRequest struct into the body of an
	// IXRTB POST request. We can do so using the net/url package.
	postBody := url.Values{}
	postBody.Set("s", dspRequest.S)
	postBody.Set("l", dspRequest.L)
	postBody.Set("d", dspRequest.D)
	postBody.Set("a", strconv.Itoa(dspRequest.A))

	// Instantiate a net/http client that can be used to make a http request to another server
	client := &http.Client{}

	// Build a new http request, which involves specifying the request type (POST), the request URL,
	// and the data that will be included in the POST request's body.
	dspRequestBody, _ := http.NewRequest("POST", dspURL, strings.NewReader(postBody.Encode()))
	// Also specify headers to specify the format of the data being sent (Content-Type, in this case
	// alphanumeric), and the length of the message.
	dspRequestBody.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	dspRequestBody.Header.Add("Content-Length", strconv.Itoa(len(postBody.Encode())))

	// 2. Send the POST Request
	log.Printf("Sending request to DSP at: %v", dspURL)

	return client.Do(dspRequestBody)
}

// parseJSONResponse is an example of how to parse a JSON response received from a server
// that responded to the initial request.
func parseJSONResponse(bid Bid, response *http.Response, dspURL string) (Bid, error) {

	// Make sure to close the connection after this method completes (defer)
	defer response.Body.Close()

	// ioutil.ReadAll uses the response.Body reader to extract the data into a byte array (body).
	// Check out this link for more about Golang's Reader and Writer interfaces:
	// https://nathanleclaire.com/blog/2014/07/19/demystifying-golangs-io-dot-reader-and-io-dot-writer-interfaces/
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error at response body from DSP at: %v", dspURL)
		return bid, err
	}

	// Using the json library we can now extract the JSON data into a bid struct
	// and populate its fields according to what was received in the response.
	err = json.Unmarshal(body, &bid)
	if err != nil {
		log.Printf("Unable to unmarshal the json body of DSP at: %v", dspURL)
		return bid, err
	}

	return bid, nil
}

// ReturnJSONResponse is an example of how to return a JSON formatted response to the entity
// that sent the request. The JSON returned will be a JSON version of the struct that was
// passed to it.
func ReturnJSONResponse(w http.ResponseWriter, result interface{}) {

	// The json Marshaller is excellent at converting Go structs into JSON objects. Simply
	// provide the struct to convert and a JSON object with the struct's key/value pairs
	// is created.
	marshalled, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(httpInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	// Using the response writer, set headers to allow cross-origin-resource-sharing and specify
	// content returned is of type JSON. Provide the httpOk 200 status code and write the marshalled
	// JSON object, which is how the response is returned.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(httpOK)
	w.Write(marshalled)
}
