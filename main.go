package main

import (
	"corona-visual-server/internal/config"
	"corona-visual-server/internal/fetcher"
	"corona-visual-server/internal/handler"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	openAPIURL = "http://openapi.data.go.kr/openapi/service/rest/Covid19/getCovid19InfStateJson"
	totalWeeks = 5
)

var netClient = &http.Client{
	Timeout: time.Second * 20,
	Transport: &http.Transport{
		DialContext:         (&net.Dialer{Timeout: 15 * time.Second}).DialContext,
		TLSHandshakeTimeout: 15 * time.Second,
	},
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	serviceKey := os.Getenv("COVID19_INFECT_SERVICE_KEY")
	if serviceKey == "" {
		log.Fatal("$SERVICE_KEY is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
		log.Println("$PORT is not set, so port set to ", port)
	}

	cfg := config.Config{
		OpenAPIURL: openAPIURL,
		ServiceKey: serviceKey,
		TotalWeeks: config.DefaultWeeks,
	}

	f := fetcher.New(&cfg, netClient)
	h := handler.New(&cfg, &f)

	r := gin.Default()
	r.GET("/", h.GetWeeklyHandler)
	r.GET("/:weeks", h.GetWeeklyHandler)
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusOK)
	})

	r.Run(":" + port)
}
