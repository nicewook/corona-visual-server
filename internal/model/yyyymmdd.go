package model

import (
	"encoding/xml"
	"time"
)

// YYYYMMDDSeoulTime 은 custom 날짜 객체로써
// YYYYMMDD 00:00 형태의 날짜를 서울 기준으로 담고 있습니다.
type YYYYMMDDSeoulTime struct {
	time.Time
}

func (c *YYYYMMDDSeoulTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	s := c.Time.Format(APIDateFormat)
	return e.EncodeElement(s, start)
}

func (c *YYYYMMDDSeoulTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var text string
	err := d.DecodeElement(&text, &start)
	if err != nil {
		return err
	}

	t, err := time.ParseInLocation(APIDateFormat, text, seoulTimeZone)
	if err != nil {
		return err
	}
	*c = YYYYMMDDSeoulTime{t}
	return nil
}
