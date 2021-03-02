package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "embed"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// func generateBarItems(data []int) []opts.BarData {
// 	items := make([]opts.BarData, len(data))
// 	for _, d := range data {
// 		items = append(items, opts.BarData{Value: d})
// 	}
// 	return items
// }

func generateBarItems() []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.BarData{Value: rand.Intn(300)})
	}
	return items
}

// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
var netTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 5 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 5 * time.Second,
}
var netClient = &http.Client{
	Timeout:   time.Second * 10,
	Transport: netTransport,
}

const (
	openAPIURL = "http://openapi.data.go.kr/openapi/service/rest/Covid19/getCovid19InfStateJson"
	dateFormat = "20060102"
)

func get3WeeksRange() (string, string) {

	cTime := time.Now()
	endDate := cTime.Format(dateFormat)
	startDate := cTime.AddDate(0, 0, -23).Format(dateFormat) // I need 21 days, but I have 23 days just in case
	fmt.Printf("startDate %v, endDate %v\n", startDate, endDate)
	return startDate, endDate
}

func getCoronaData() ([]byte, error) {

	// make request with query https://stackoverflow.com/a/30657518/6513756
	req, err := http.NewRequest("GET", openAPIURL, nil)
	if err != nil {
		return nil, err
	}

	startDate, endDate := get3WeeksRange()
	q := req.URL.Query()
	q.Add("serviceKey", serviceKey)
	q.Add("pageNo", "1")
	q.Add("numOfRows", "25") // I will have max 23 days result
	q.Add("startCreateDt", startDate)
	q.Add("endCreateDt", endDate)

	fmt.Println("req.URL.String():", req.URL.String())
	req.URL.RawQuery = q.Encode() // this make added query to attached
	fmt.Println("req.URL.String():", req.URL.String())

	resp, err := netClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// response
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}

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

func getAddCount(today Item, yday Item) string {

	tCareCnt, err := strconv.Atoi(today.CareCnt)
	if err != nil {
		return "-1"
	}
	yCareCnt, err := strconv.Atoi(yday.CareCnt)
	if err != nil {
		return "-1"
	}
	tClearCnt, err := strconv.Atoi(today.ClearCnt)
	if err != nil {
		return "-1"
	}
	yClearCnt, err := strconv.Atoi(yday.ClearCnt)
	if err != nil {
		return "-1"
	}
	tDeathCnt, err := strconv.Atoi(today.DeathCnt)
	if err != nil {
		return "-1"
	}
	yDeathCnt, err := strconv.Atoi(yday.DeathCnt)
	if err != nil {
		return "-1"
	}

	return strconv.Itoa(tCareCnt + tClearCnt + tDeathCnt - yCareCnt - yClearCnt - yDeathCnt)
}

func weeklyHandler(w http.ResponseWriter, r *http.Request) {

	// if the last creation of the html is over 2 min
	// - (24*60*60)/1000 = 86.4 seconds // 1000 request per day
	b, err := getCoronaData()
	if err != nil {
		log.Println(err)
		// need return code
		return
	}
	// fmt.Println(string(b))

	var resp Response
	if err := xml.Unmarshal(b, &resp); err != nil {
		log.Println(err)
		return
	}

	var data []CoronaDailyData
	for i := range resp.Body.Items.Item {

		if i == len(resp.Body.Items.Item)-1 {
			continue
		}
		t, err := time.Parse(dateFormat, resp.Body.Items.Item[i].StateDt)
		if err != nil {
			log.Println(err)
			continue
		}

		var d CoronaDailyData
		d.Date = t.AddDate(0, 0, -1).Format(dateFormat)
		d.AddCount = getAddCount(resp.Body.Items.Item[i], resp.Body.Items.Item[i+1])
		data = append(data, d)
	}

	// fmt.Printf("data: %+v\n", data)

	// create a new bar instance
	bar := charts.NewBar()
	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "My first bar chart generated by go-echarts",
			Subtitle: "It's extremely easy to use, right?",
		}),
		charts.WithLegendOpts(opts.Legend{Show: true}),
	)

	// Put data into instance
	bar.SetXAxis([]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}).
		AddSeries("Category A", generateBarItems()).
		AddSeries("Category B", generateBarItems()).
		AddSeries("Category C", generateBarItems())
	// Where the magic happens
	f, _ := os.Create("bar.html")
	bar.Render(f)

	htmlFile := "./bar.html"
	http.ServeFile(w, r, htmlFile)
}

const port = ":8081"

//go:embed coronaState.key
var serviceKey string

func main() {

	fmt.Println("service Key: ", serviceKey)
	http.HandleFunc("/", weeklyHandler)
	http.ListenAndServe(port, nil)
}
