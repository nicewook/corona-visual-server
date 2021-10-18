package handler

import (
	"corona-visual-server/internal/config"
	"corona-visual-server/internal/fetcher"
	"corona-visual-server/internal/model"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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

type URITotalWeeks struct {
	Weeks string `uri:"weeks"`
}

// GetWeeklyHandler handles weekly request.
func (h *Handler) GetWeeklyHandler(c *gin.Context) {

	h.config.TotalWeeks = config.DefaultWeeks

	// get totalWeeks from uri, if exist
	var tw URITotalWeeks
	if err := c.ShouldBindUri(&tw); err != nil {
		log.Println("err: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
	} else if tw.Weeks != "" {
		if weeks, err := strconv.Atoi(tw.Weeks); err != nil {
			log.Println("err: ", err)
			c.AbortWithStatus(http.StatusBadRequest)
		} else {
			log.Println("total weeks to display: ", weeks)
			h.config.TotalWeeks = weeks
		}
	}

	coronaDataSet, err := h.fetcher.GetCoronaData()
	if err != nil {
		log.Printf("h.fetcher.GetCoronaData() returns an err = %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	data := coronaDataSet.Data

	// set some global options like Title/Legend/ToolTip or anything else
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Covid confirmed person data comparison",
			Subtitle: fmt.Sprintf("%d Weeks comparison of each weekday", h.config.TotalWeeks),
			Left:     "5%",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: true,
			Left: "48%",
			Top:  "5%",
		}),
		charts.WithInitializationOpts(opts.Initialization{
			PageTitle: "Weekly comparison of the Corona confirmed persion in South Korea", // HTML title
			Width:     "950px",                                                            // Width of canvas
			Height:    "550px",                                                            // Height of canvas
		}),
	)

	// Put data into the bar instance

	for i := 0; i < h.config.TotalWeeks; i++ {
		label := fmt.Sprintf("%d weeks ago", h.config.TotalWeeks-i)
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

	if err := bar.Render(c.Writer); err != nil {
		log.Println(err)
	}
}

// getWeeklyAxis finds the starting weekday of the xAxis
func (h *Handler) getWeeklyAxis(data model.CoronaDailyData) []string {
	dateFormat := "2006-01-02"
	wDay := data.Date.Weekday().String()
	log.Printf("starting date: %v, starting date local: %v", data.Date.Format(dateFormat), data.Date.Local().Format(dateFormat))
	log.Printf("weekday start: %v ", wDay)

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
