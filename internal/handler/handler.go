package handler

import (
	"corona-visual-server/internal/config"
	"corona-visual-server/internal/fetcher"
	"corona-visual-server/internal/model"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var weekdays = []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}

// Handler is the main handler.
type Handler struct {
	fetcher *fetcher.Fetcher
	config  *config.Config
}

// New returns the Handler.
func New(config *config.Config, fetcher *fetcher.Fetcher) Handler {
	return Handler{
		fetcher: fetcher,
		config:  config,
	}
}

// GetWeeklyHandler handles weekly request.
func (h *Handler) GetWeeklyHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetWeeklyHandler request = %+v", r)
	if r.Method != http.MethodGet || r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	coronaDataSet, err := h.fetcher.GetCoronaData()
	if err != nil {
		log.Printf("h.fetcher.GetCoronaData() returns an err = %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data := coronaDataSet.Data

	// create a new bar instance
	bar := charts.NewBar()
	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Covid confirmed person data comparison",
			Subtitle: fmt.Sprintf("%d Weeks comparison of each weekday", h.config.TotalWeeks),
			Left:     "5%",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: true,
			Left: "48%",
			Top:  "7%",
		}),
	)

	// Put data into the bar instance

	for i := 0; i < h.config.TotalWeeks; i++ {
		label := fmt.Sprintf("%d weeks ago", 5-i)
		start := i * 7
		end := start + 7
		bar.AddSeries(label, generateWeeklyBarItems(data[start:end]))
	}
	bar.SetXAxis(h.getWeeklyAxis(data[0])).
		SetSeriesOptions(charts.WithLabelOpts(opts.Label{
			Show:     true,
			Position: "top",
		}),
		)

	if err := bar.Render(w); err != nil {
		log.Println(err)
	}
}

// getWeeklyAxis finds the starting weekday of the xAxis
func (h *Handler) getWeeklyAxis(data model.CoronaDailyData) []string {
	wDay := data.Date.Weekday().String()
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
