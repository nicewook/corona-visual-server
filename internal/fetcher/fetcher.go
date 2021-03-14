package fetcher

import (
	"corona-visual-server/internal/config"
	"corona-visual-server/internal/model"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
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
	log.Printf("startDate=%v, endDate=%v", startDate, endDate)
	return startDate, endDate
}

// GetCoronaData returns CoronaData
func (f *Fetcher) GetCoronaData() (*model.CoronaDailyDataResult, error) {
	// make request with query https://stackoverflow.com/a/30657518/6513756
	modelResponse, err := RequestCoronaAPI(f)
	if err != nil {
		return nil, err
	}

	return ProcessResponse(modelResponse, f.config.DateFormat)
}

// ResultCoronaAPI는 보건복지부 API로 요청을 보내 지난 25일간 코로나 일일 현황을 받아옵니다.
func RequestCoronaAPI(f *Fetcher) (*model.Response, error) {
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

	var modelResponse model.Response
	if err := xml.Unmarshal(bodyBytes, &modelResponse); err != nil {
		return nil, err
	}
	return &modelResponse, nil
}

// ProcessResponse 는 코로나API에서 받은 XML 데이터를 프로세스해 차트를 위한 데이터로 변환합니다.
func ProcessResponse(modelResponse *model.Response, dateFormat string) (*model.CoronaDailyDataResult, error) {
	var data []model.CoronaDailyData
	for i, todayCoronaData := range modelResponse.Body.Items.Item {
		if i == len(modelResponse.Body.Items.Item)-1 {
			continue
		}
		t, err := time.ParseInLocation(dateFormat, todayCoronaData.StateDt, config.SeoulTZ)
		if err != nil {
			log.Println(err)
			continue
		}

		var d model.CoronaDailyData

		// TODO: Explain why this needs to be subtracted.
		d.Date = t.AddDate(0, 0, -1).Format(dateFormat)
		d.AddCount = getAddCount(todayCoronaData, modelResponse.Body.Items.Item[i+1])
		data = append(data, d)
	}

	// 21일 간 데이터 가져오기
	if data == nil || len(data) < 21 {
		return nil, fmt.Errorf("expected len(data) < 21, but received len(data) = %v, data = %v", len(data), data)
	}
	data = data[:21]

	// reverse and get exact 21 data
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}

	return &model.CoronaDailyDataResult{Data: data}, nil
}

func getAddCount(today model.Item, yday model.Item) int64 {
	return today.CareCnt + today.ClearCnt + today.DeathCnt - yday.CareCnt - yday.ClearCnt - yday.DeathCnt
}
