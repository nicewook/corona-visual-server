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

// Item represents an individual item of Response.
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

// CoronaDailyData is a single data point.
type CoronaDailyData struct {
	Date     string
	AddCount string
}
