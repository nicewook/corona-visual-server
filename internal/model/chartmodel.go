package model

// CoronaDailyData is a single data point.
type CoronaDailyData struct {
	Date     YYYYMMDDSeoulTime `xml:"date"`
	AddCount int64             `xml:"addCount"`
}

// CoronaDailyDataResult 는 CoronaDailyData 들을 담고 있는 객체입니다.
// XML 에서는 Root Node가 꼭 하나는 필요하기 때문에 Marshal에 사용됩니다.
type CoronaDailyDataResult struct {
	Data []CoronaDailyData `xml:"data"`
}
