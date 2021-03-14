package model

import "encoding/xml"

// Response represents the result of research.swtch.com
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
			Item []Item `xml:"item"`
		} `xml:"items"`
		NumOfRows  string `xml:"numOfRows"`
		PageNo     string `xml:"pageNo"`
		TotalCount string `xml:"totalCount"`
	} `xml:"body"`
}

// Item represents an individual item of Response.
type Item struct {
	Text           string `xml:",chardata"`
	AccDefRate     string `xml:"accDefRate"`
	AccExamCnt     string `xml:"accExamCnt"`
	AccExamCompCnt string `xml:"accExamCompCnt"`
	CareCnt        int64  `xml:"careCnt"`
	ClearCnt       int64  `xml:"clearCnt"`
	CreateDt       string `xml:"createDt"`
	DeathCnt       int64  `xml:"deathCnt"`
	DecideCnt      int64  `xml:"decideCnt"`
	ExamCnt        int64  `xml:"examCnt"`
	// 오타 실화입니까 보건복지부 ㅜㅜ..
	ResutlNegCnt int64  `xml:"resutlNegCnt"`
	Seq          int64  `xml:"seq"`
	StateDt      string `xml:"stateDt"`
	StateTime    string `xml:"stateTime"`
	UpdateDt     string `xml:"updateDt"`
}

// CoronaDailyData is a single data point.
type CoronaDailyData struct {
	Date     string `xml:"date"`
	AddCount int64  `xml:"addCount"`
}

type CoronaDailyDataResult struct {
	Data []CoronaDailyData `xml:"data"`
}
