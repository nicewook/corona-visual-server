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
)

const (
	openAPIURL = "http://openapi.data.go.kr/openapi/service/rest/Covid19/getCovid19InfStateJson"
)

var netClient = &http.Client{
	Timeout: time.Second * 20,
	Transport: &http.Transport{
		DialContext:         net.Dialer{Timeout: 15 * time.Second}.DialContext,
		TLSHandshakeTimeout: 15 * time.Second,
	},
}

func main() {
	serviceKey := os.Getenv("SERVICE_KEY")
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
	}

	f := fetcher.New(&cfg, netClient)
	h := handler.New(&cfg, &f)

	http.HandleFunc("/", h.GetWeeklyHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
