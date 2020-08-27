package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"sync"
	"testing"
	"ums/dbops"
	"ums/platform"
)

// initial setup
func setupForTest() (*gorm.DB, map[string]chan bool, *sync.RWMutex, *httptest.Server) {
	db, errDB := gorm.Open("mysql", "root:Anish@6030@tcp(127.0.0.1:3306)/test_database?charset=utf8&parseTime=True")

	if errDB != nil {
		fmt.Printf("Expected no error, got %v\n", errDB)
		os.Exit(1)
	}
	db.DropTableIfExists(&dbops.DataInDB{})
	db.Debug().AutoMigrate(&dbops.DataInDB{})

	// ChannelMap variable to hold map of channels
	channelMap := make(map[string]chan bool)

	// Lock variable
	var lock = sync.RWMutex{}
	////Grab our router
	router := SetupRouter(db, channelMap, &lock)
	ts := httptest.NewServer(router)
	return db, channelMap, &lock, ts

}

func testForCreateURLEndpoint(ts *httptest.Server, channelMap map[string]chan bool, lock *sync.RWMutex) (*http.Response, error) {
	form := url.Values{}
	form.Add("crawl_timeout", "3")
	form.Add("frequency", "3")
	form.Add("failure_threshold", "50")
	form.Add("url", "https://www.khjfcxdoooogle.com/")

	resp, err := http.PostForm(fmt.Sprintf("%s/url", ts.URL), form)
	stopGoRoutine("https://www.khjfcxdoooogle.com/", channelMap, lock)
	return resp, err
}

func parseJSON(resp *http.Response) (map[string]interface{}, error) {
	var respJSON map[string]interface{}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	errJSON := json.Unmarshal(bodyBytes, &respJSON)
	return respJSON, errJSON

}

func stopGoRoutine(crawlURL string, channelMap map[string]chan bool, lock *sync.RWMutex) {
	channel, exists := platform.ReadChanFromChanMap(crawlURL, channelMap, lock)
	if exists {
		channel <- true
	}
}

// TestCreateURLEntry function to test the CreateURLEntry endpoint
func TestEndpoints(t *testing.T) {
	db, channelMap, lock, ts := setupForTest()
	defer db.Close()
	defer ts.Close()

	// TestCreateURL Endpoint
	resp, err := testForCreateURLEndpoint(ts, channelMap, lock) //http.PostForm(fmt.Sprintf("%s/url", ts.URL), form)
	assert.Nilf(t, err, "Expected no error, got %v", err)
	assert.EqualValuesf(t, 200, resp.StatusCode, "Expected status code 200, got %v", resp.StatusCode)
	jsonRespCreateURL, jsonErrCreateURL := parseJSON(resp)
	assert.Nilf(t, jsonErrCreateURL, "Expected no error, got %v", err)
	_, exists := jsonRespCreateURL["id"]
	assert.True(t, exists)

}
