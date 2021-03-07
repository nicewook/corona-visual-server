package fetcher

import (
	"corona-visual-server/internal/config"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Fetcher fetches Corona data.
type Fetcher struct {
	config *config.Config
	client *http.Client
}

// New returns Fetcher
func New(config *config.Config, client *http.Client) Fetcher {
	return Fetcher{
		config: config,
		client: client,
	}
}

func (f *Fetcher) get3WeeksRange() (string, string) {
	cTime := time.Now()
	endDate := cTime.Format(f.config.DateFormat)
	startDate := cTime.AddDate(0, 0, -23).Format(f.config.DateFormat) // I need 21 days, but I have 23 days just in case
	fmt.Printf("startDate %v, endDate %v\n", startDate, endDate)
	return startDate, endDate
}

// GetCoronaData returns CoronaData
// TODO: This function should return CoronaStruct instead of []byte
func (f *Fetcher) GetCoronaData() ([]byte, error) {
	// make request with query https://stackoverflow.com/a/30657518/6513756
	fmt.Println("GetCoronaData")

	req, err := http.NewRequest("GET", f.config.OpenAPIURL, nil)
	if err != nil {
		return nil, err
	}

	startDate, endDate := f.get3WeeksRange()
	q := req.URL.Query()
	q.Add("serviceKey", f.config.ServiceKey)
	q.Add("pageNo", "1")
	q.Add("numOfRows", "25") // I will have max 23 days result
	q.Add("startCreateDt", startDate)
	q.Add("endCreateDt", endDate)

	req.URL.RawQuery = q.Encode() // this make added query to attached AND URL encoding
	// fmt.Println("req.URL.String():", req.URL.String())

	resp, err := f.client.Do(req)
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
	fmt.Println("GetCoronaData success")

	return bodyBytes, nil
}
