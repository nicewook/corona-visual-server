package model

import (
	"log"
	"time"
)

// APIDateFormat 은 보건복지부 API 에서 사용 되는 날짜 포맷입니다.
// 예를 들어, Response 중 StateDt 필드의 날짜 값 형식에 사용됩니다.
const APIDateFormat = "20060102"

// 보건복지부 날짜는 모두 Seoul 기준입니다.
var seoulTimeZone *time.Location

func init() {
	var err error
	seoulTimeZone, err = time.LoadLocation("Asia/Seoul")
	if err != nil {
		log.Fatalf("failed to load the timezone Asia/Seoul. Please update the local time database")
	}
}
