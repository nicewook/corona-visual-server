package fetcher

import (
	"corona-visual-server/internal/model"
	"encoding/xml"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

func TestProcessResponse(t *testing.T) {
	sampleResponse := getSampleResponse()

	response, err := ProcessResponse(sampleResponse)
	if err != nil {
		t.Errorf("ProcessResponse should not have failed but failed. err = %v", err)
	}

	expected := getExpectedOutput()
	if !reflect.DeepEqual(response, expected) {
		t.Errorf("got = %v\nexpected = %v", response, expected)
	}
}

func getSampleResponse() *model.Response {
	return mustConvertFileToResponse(
		"../testdata/sample_response.xml",
		&model.Response{}).(*model.Response)
}

func getExpectedOutput() *model.CoronaDailyDataResult {
	return mustConvertFileToResponse(
		"../testdata/expected_processed_output.xml",
		&model.CoronaDailyDataResult{}).(*model.CoronaDailyDataResult)
}

func mustConvertFileToResponse(filepath string, response interface{}) interface{} {
	bytes := mustReadFile(filepath)
	return mustConvert(bytes, response)
}

func mustReadFile(filepath string) []byte {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalf("failed to read %v. maybe the file does not exist? err = %v", filepath, err)
	}
	return bytes
}

func mustConvert(b []byte, obj interface{}) interface{} {
	if err := xml.Unmarshal(b, obj); err != nil {
		log.Fatalf("failed to unmarshal. Please make sure the file is a correct XML. err = %v", err)
	}
	return obj
}
