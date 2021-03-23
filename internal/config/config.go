package config

// Config represents Application global config.
type Config struct {
	// 보건복지부 API URL 을 의미합니다.
	OpenAPIURL string
	// 보건복지부 API 의 서비스 키를 의미합니다.
	ServiceKey string
	// 몇 주간의 확진자 수를 비교할지 결정한다
	TotalWeeks int
}
