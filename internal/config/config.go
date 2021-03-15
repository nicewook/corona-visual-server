package config

import (
	"log"
	"time"
)

// Config represents Application global config.
type Config struct {
	OpenAPIURL string
	DateFormat string
	ServiceKey string
	Port       string
}

var SeoulTZ *time.Location

func init() {
	var err error
	SeoulTZ, err = time.LoadLocation("Asia/Seoul")
	if err != nil {
		log.Fatalf("failed to load the timezone Asia/Seoul. Please update the local time database")
	}
}
