package main

import (
	_ "embed"
	"encoding/xml"
	"net"
	"net/http"
	"time"
)

// types and constants and variables

// research.swtch.com
type Response struct {
	XMLName xml.Name `xml:"response"`
	Text    string   `xml:",chardata"`
	Header  struct {
		Text       string `xml:",chardata"`
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	Body struct {
		Text  string `xml:",chardata"`
		Items struct {
			Text string `xml:",chardata"`
			Item []struct {
				Text           string `xml:",chardata"`
				AccDefRate     string `xml:"accDefRate"`
				AccExamCnt     string `xml:"accExamCnt"`
				AccExamCompCnt string `xml:"accExamCompCnt"`
				CareCnt        string `xml:"careCnt"`
				ClearCnt       string `xml:"clearCnt"`
				CreateDt       string `xml:"createDt"`
				DeathCnt       string `xml:"deathCnt"`
				DecideCnt      string `xml:"decideCnt"`
				ExamCnt        string `xml:"examCnt"`
				ResutlNegCnt   string `xml:"resutlNegCnt"`
				Seq            string `xml:"seq"`
				StateDt        string `xml:"stateDt"`
				StateTime      string `xml:"stateTime"`
				UpdateDt       string `xml:"updateDt"`
			} `xml:"item"`
		} `xml:"items"`
		NumOfRows  string `xml:"numOfRows"`
		PageNo     string `xml:"pageNo"`
		TotalCount string `xml:"totalCount"`
	} `xml:"body"`
}

type Item struct {
	Text           string `xml:",chardata"`
	AccDefRate     string `xml:"accDefRate"`
	AccExamCnt     string `xml:"accExamCnt"`
	AccExamCompCnt string `xml:"accExamCompCnt"`
	CareCnt        string `xml:"careCnt"`
	ClearCnt       string `xml:"clearCnt"`
	CreateDt       string `xml:"createDt"`
	DeathCnt       string `xml:"deathCnt"`
	DecideCnt      string `xml:"decideCnt"`
	ExamCnt        string `xml:"examCnt"`
	ResutlNegCnt   string `xml:"resutlNegCnt"`
	Seq            string `xml:"seq"`
	StateDt        string `xml:"stateDt"`
	StateTime      string `xml:"stateTime"`
	UpdateDt       string `xml:"updateDt"`
}

type CoronaDailyData struct {
	Date     string
	AddCount string
}

// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
var netTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 15 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 15 * time.Second,
}

var netClient = &http.Client{
	Timeout:   time.Second * 20,
	Transport: netTransport,
}

var (
	//go:embed service.key
	serviceKey string
	weekdays   = []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
)

const (
	openAPIURL = "http://openapi.data.go.kr/openapi/service/rest/Covid19/getCovid19InfStateJson"
	dateFormat = "20060102"
	port       = ":8081"
)
