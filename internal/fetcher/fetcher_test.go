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

	response, err := ProcessResponse(sampleResponse, "20060102")
	if err != nil {
		t.Errorf("ProcessResponse should not have failed but failed. err = %v", err)
	}

	expected := getExpectedOutput()
	if !reflect.DeepEqual(response, expected) {
		t.Errorf("got = %v\nexpected = %v", response, expected)
	}
}

func getExpectedOutput() *model.CoronaDailyDataResult {
	filepath := "../testdata/expected_processed_output.xml"
	bytes := mustReadFile(filepath)
	response := &model.CoronaDailyDataResult{}
	if err := xml.Unmarshal(bytes, response); err != nil {
		log.Fatalf("failed to unmarshal to model.CoronaDailyDataResult. Please make sure the file is a correct XML. err = %v", err)
	}
	return response
}

func getSampleResponse() *model.Response {
	filepath := "../testdata/sample_response.xml"
	bytes := mustReadFile(filepath)
	response := &model.Response{}
	if err := xml.Unmarshal(bytes, response); err != nil {
		log.Fatalf("failed to unmarshal to model.Response. Please make sure the file is a correct XML. err = %v", err)
	}
	return response
}

func mustReadFile(filepath string) []byte {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalf("failed to read %v. maybe the file does not exist? err = %v", filepath, err)
	}
	return bytes
}
