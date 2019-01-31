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

	client := &http.Client{}
	dspRequestBody, _ := http.NewRequest("POST", hObj.Config.Dspurl, strings.NewReader(data.Encode()))

	dspRequestBody.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	dspRequestBody.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	hObj.Logger.Info("sending request to dsp: ", hObj.Config.Dspurl)
	response, err := client.Do(dspRequestBody)
	if err != nil {
		hObj.Logger.Error("Error at sending request to DSP: ", err)
		HandleSuccess(&w, HeartBeatResponse{Okay: false}, hObj.Logger)
		return
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		hObj.Logger.Error("Error at response body:", err)
		HandleSuccess(&w, HeartBeatResponse{Okay: false}, hObj.Logger)
		return
	}

	var bid Bid
	err = json.Unmarshal(body, &bid)
	if err != nil {
		hObj.Logger.Error("Unable to unmarshal the json body:", err)
		HandleSuccess(&w, HeartBeatResponse{Okay: false}, hObj.Logger)
		return
	}

	HandleSuccess(&w, bid, hObj.Logger)
}

func (hObj *HandlerObj) Heartbeat(w http.ResponseWriter, r *http.Request) {
	hObj.Logger.Info("/heartbeat request received")
	HandleSuccess(&w, HeartBeatResponse{Okay: true}, hObj.Logger)
}
