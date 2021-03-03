package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func generateWeeklyBarItems(data []CoronaDailyData) []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.BarData{Value: data[i].AddCount})
	}
	return items
}

func get3WeeksRange() (string, string) {
	cTime := time.Now()
	endDate := cTime.Format(dateFormat)
	startDate := cTime.AddDate(0, 0, -23).Format(dateFormat) // I need 21 days, but I have 23 days just in case
	fmt.Printf("startDate %v, endDate %v\n", startDate, endDate)
	return startDate, endDate
}

func getCoronaData() ([]byte, error) {
	// make request with query https://stackoverflow.com/a/30657518/6513756
	fmt.Println("getCoronaData")

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

	req.URL.RawQuery = q.Encode() // this make added query to attached AND URL encoding
	// fmt.Println("req.URL.String():", req.URL.String())

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
	fmt.Println("getCoronaData success")

	return bodyBytes, nil
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

// getWeeklyAxis finds the starting weekday of the xAxis
func getWeeklyAxis(data CoronaDailyData) []string {
	t, err := time.Parse(dateFormat, data.Date)
	if err != nil {
		log.Println(err)
		return weekdays
	}
	wDay := t.Weekday().String()
	fmt.Println("weekday start: ", wDay)

	var idx int
	for i, d := range weekdays {
		if strings.Contains(wDay, d) {
			idx = i
		}
	}
	result := append(weekdays[idx:], weekdays[:idx]...)
	return result
}

func weeklyHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("weeklyHandler")
	b, err := getCoronaData()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// fmt.Println(string(b))

	var resp Response
	if err := xml.Unmarshal(b, &resp); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
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

	// reverse and get exact 21 data
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	cutCount := len(data) - 21 // 3 weeks == 21 days
	data = data[cutCount:]

	// create a new bar instance
	bar := charts.NewBar()
	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Covid confirmed person data comparison",
			Subtitle: "3 Weeks comparison of each weekday",
			Left:     "5%",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: true,
			Left: "48%",
			Top:  "7%",
		}),
	)

	// Put data into the bar instance
	bar.SetXAxis(getWeeklyAxis(data[0])).
		AddSeries("3rd weeks ago", generateWeeklyBarItems(data[:7])).
		AddSeries("2nd weeks ago", generateWeeklyBarItems(data[7:14])).
		AddSeries("1st weeks ago", generateWeeklyBarItems(data[14:])).
		SetSeriesOptions(charts.WithLabelOpts(opts.Label{
			Show:     true,
			Position: "top",
		}),
		)
	bar.Render(w)
	fmt.Println("--")
}

func main() {

	serviceKey = os.Getenv("SERVICE_KEY")
	if serviceKey == "" {
		log.Fatal("$SERVICE_KEY is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
		log.Println("$PORT is not set, so port set to ", port)
	}

	// log.Printf("service key: %v, port %v\n", serviceKey, port)

	http.HandleFunc("/", weeklyHandler)
	http.ListenAndServe(":"+port, nil)
}
