package fetcher

import (
	"corona-visual-server/internal/config"
	"corona-visual-server/internal/model"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

func (f *Fetcher) getWeeksRange() (string, string) {
	cTime := time.Now()
	endDate := cTime.Format(model.APIDateFormat)
	startDate := cTime.AddDate(0, 0, -(f.config.TotalWeeks*7 + 2)).Format(model.APIDateFormat) // I need 21 days, but I have 23 days just in case
	log.Printf("startDate=%v, endDate=%v", startDate, endDate)
	return startDate, endDate
}

// GetCoronaData returns CoronaData
func (f *Fetcher) GetCoronaData() (*model.CoronaDailyDataResult, error) {
	// make request with query https://stackoverflow.com/a/30657518/6513756
	modelResponse, err := RequestCoronaAPI(f, f.config.TotalWeeks)
	if err != nil {
		return nil, err
	}

	return f.ResponseToChartData(modelResponse)
}

// ResultCoronaAPI는 보건복지부 API로 요청을 보내 지난 25일간 코로나 일일 현황을 받아옵니다.
func RequestCoronaAPI(f *Fetcher, weeks int) (*model.Response, error) {
	req, err := http.NewRequest("GET", f.config.OpenAPIURL, nil)
	if err != nil {
		return nil, err
	}

	startDate, endDate := f.getWeeksRange()
	numOfRows := weeks*7 + 2
	q := req.URL.Query()
	q.Add("serviceKey", f.config.ServiceKey)
	q.Add("pageNo", "1")
	q.Add("numOfRows", strconv.Itoa(numOfRows)) // I will have max 23 days result
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

// ResponseToChartData 는 코로나API에서 받은 XML 데이터를 프로세스해 차트를 위한 데이터로 변환합니다.
func (f *Fetcher) ResponseToChartData(modelResponse *model.Response) (*model.CoronaDailyDataResult, error) {
	var data []model.CoronaDailyData
	for i, todayCoronaData := range modelResponse.Body.Items.Item {
		if i == len(modelResponse.Body.Items.Item)-1 {
			continue
		}

		var d model.CoronaDailyData

		/* XML 의 Item 참고
		<item>
			<accDefRate>1.3637663657</accDefRate>
			<accExamCnt>7046782</accExamCnt>
			<accExamCompCnt>6978908</accExamCompCnt>
			<careCnt>6884</careCnt>
			<clearCnt>86625</clearCnt>
			<createDt>2021-03-13 09:36:41.886</createDt>
			<deathCnt>1667</deathCnt>
			<decideCnt>95176</decideCnt>
			<examCnt>67874</examCnt>
			<resutlNegCnt>6883732</resutlNegCnt>
			<seq>447</seq>
			<stateDt>20210313</stateDt>
			<stateTime>00:00</stateTime>
			<updateDt>null</updateDt>
		</item>
		*/
		// StateDt와 Time을 보면 집계를 한 바로 다음날인 자정 00:00시 인 것을 알 수 있다.
		// 따라서 d.Date는 하루 전이어야 한다.
		d.Date = model.YYYYMMDDSeoulTime{Time: todayCoronaData.StateDt.Time.AddDate(0, 0, -1)}
		d.AddCount = getAddCount(todayCoronaData, modelResponse.Body.Items.Item[i+1])
		data = append(data, d)
	}

	// 원하는 주 만큼을 가져오기
	if data == nil || len(data) < f.config.TotalWeeks*7 {
		return nil, fmt.Errorf("expected len(data) < 21, but received len(data) = %v, data = %v", len(data), data)
	}
	data = data[:f.config.TotalWeeks*7]

	// reverse and get exact 21 data
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}

	return &model.CoronaDailyDataResult{Data: data}, nil
}

func getAddCount(today model.Item, yday model.Item) int64 {
	return today.CareCnt + today.ClearCnt + today.DeathCnt - yday.CareCnt - yday.ClearCnt - yday.DeathCnt
}
