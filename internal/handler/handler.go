package handler

import (
	"corona-visual-server/internal/config"
	"corona-visual-server/internal/fetcher"
	"corona-visual-server/internal/model"
	"encoding/xml"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var weekdays = []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}

// Handler is the main handler.
type Handler struct {
	fetcher *fetcher.Fetcher
	config  *config.Config
}

// New returns the Handler.
func New(fetcher *fetcher.Fetcher) Handler {
	return Handler{
		fetcher: fetcher,
	}
}

// GetWeeklyHandler handles weekly request.
// TODO: This function needs refactoring.
func (h *Handler) GetWeeklyHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("weeklyHandler")
	b, err := h.fetcher.GetCoronaData()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// fmt.Println(string(b))
	var resp model.Response
	if err := xml.Unmarshal(b, &resp); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var data []model.CoronaDailyData
	for i := range resp.Body.Items.Item {
		if i == len(resp.Body.Items.Item)-1 {
			continue
		}
		t, err := time.Parse(h.config.DateFormat, resp.Body.Items.Item[i].StateDt)
		if err != nil {
			log.Println(err)
			continue
		}

		var d model.CoronaDailyData
		d.Date = t.AddDate(0, 0, -1).Format(h.config.DateFormat)
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
	bar.SetXAxis(h.getWeeklyAxis(data[0])).
		AddSeries("3rd weeks ago", generateWeeklyBarItems(data[:7])).
		AddSeries("2nd weeks ago", generateWeeklyBarItems(data[7:14])).
		AddSeries("1st weeks ago", generateWeeklyBarItems(data[14:])).
		SetSeriesOptions(charts.WithLabelOpts(opts.Label{
			Show:     true,
			Position: "top",
		}),
		)
	err = bar.Render(w)
	if err != nil {
		log.Println(err)
	}
}

func getAddCount(today model.Item, yday model.Item) string {

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
func (h *Handler) getWeeklyAxis(data model.CoronaDailyData) []string {
	t, err := time.Parse(h.config.DateFormat, data.Date)
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

func generateWeeklyBarItems(data []model.CoronaDailyData) []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.BarData{Value: data[i].AddCount})
	}
	return items
}
