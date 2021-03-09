package main

import (
	"net"
	"net/http"
	"time"
)

// types and constants and variables

// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
var netTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 15 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 15 * time.Second,
}

var netClient = &http.Client{
	Timeout:   time.Second * 20,
	Transport: netTransport,
}

const (
	openAPIURL = "http://openapi.data.go.kr/openapi/service/rest/Covid19/getCovid19InfStateJson"
	dateFormat = "20060102"
)
