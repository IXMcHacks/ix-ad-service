package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/eftakhairul/ix-ad-service/lib"
	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"
)

//Handler context obj
type HandlerObj struct {
	Logger *logrus.Logger
	Config *lib.Config
}
type Bid struct {
	Advertiser    string `json:"Advertiser"`
	BidPrice      int    `json:"BidPrice"`
	AdURL         string `json:"AdURL"`
	AdDescription string `json:"AdDescription"`
}

type DspResponse struct {
	Bids []Bid `json:"bids"`
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

func (hObj *HandlerObj) Adserving(w http.ResponseWriter, r *http.Request) {
	//```10.4.145.10/ixrtb?s=wide&l=cooking&d=working&a=33&code=yes```
	var decoder = schema.NewDecoder()
	var dspRequest DspRequest
	hObj.Logger.Info("Ad request received")

	err := decoder.Decode(&dspRequest, r.URL.Query())
	if err != nil {
		hObj.Logger.Error(fmt.Sprintf("Error in GET parameters : %v", err))
		HandleSuccess(&w, HeartBeatResponse{Okay: false}, hObj.Logger)
		return
	}

	data := url.Values{}
	data.Set("s", dspRequest.S)
	data.Set("l", dspRequest.L)
	data.Set("d", dspRequest.D)
	data.Set("a", strconv.Itoa(dspRequest.A))

	bid1, err := sendDSP(data, hObj.Config.DspUrlOne, hObj.Logger)
	if err != nil {
		hObj.Logger.Error("error", err)
		HandleSuccess(&w, HeartBeatResponse{Okay: false}, hObj.Logger)
		return
	}

	bid2, err := sendDSP(data, hObj.Config.DspUrlTwo, hObj.Logger)
	if err != nil {
		hObj.Logger.Error("error", err)
		HandleSuccess(&w, HeartBeatResponse{Okay: false}, hObj.Logger)
		return
	}

	if bid1.BidPrice > bid2.BidPrice {
		HandleSuccess(&w, bid1, hObj.Logger)
		return
	}

	HandleSuccess(&w, bid2, hObj.Logger)
}

func (hObj *HandlerObj) Heartbeat(w http.ResponseWriter, r *http.Request) {
	hObj.Logger.Info("/heartbeat request received")
	HandleSuccess(&w, HeartBeatResponse{Okay: true}, hObj.Logger)
}

func sendDSP(data url.Values, dspUrl string, logger *logrus.Logger) (Bid, error) {
	var bid Bid
	client := &http.Client{}

	dspRequestBody, _ := http.NewRequest("POST", dspUrl, strings.NewReader(data.Encode()))

	dspRequestBody.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	dspRequestBody.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	logger.Info("sending request to dsp: ", dspUrl)
	response, err := client.Do(dspRequestBody)

	if err != nil {
		logger.Error("Error at sending request to DSP: ", dspUrl)
		return bid, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error("Error at response body from DSP:", dspUrl)
		return bid, err
	}

	err = json.Unmarshal(body, &bid)
	if err != nil {
		logger.Error("Unable to unmarshal the json body of DSP:", dspUrl)
		return bid, err
	}

	return bid, nil
}
