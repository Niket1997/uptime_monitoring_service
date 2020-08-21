package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"ums/dbops"
	"ums/platform"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func performRequest(r http.Handler, method, path string, form url.Values) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, strings.NewReader(form.Encode()))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
func TestCreateURLEntry(t *testing.T) {
	// Build our expected body
	db, err := gorm.Open("mysql", "root:Anish@6030@tcp(127.0.0.1:3306)/test_database?charset=utf8&parseTime=True")

	assert.Nil(t, err)
	defer db.Close()

	db.DropTableIfExists(&dbops.DataInDB{})
	db.Debug().AutoMigrate(&dbops.DataInDB{})

	// ChannelMap variable to hold map of channels
	channelMap := make(map[string]chan bool)

	// Lock variable
	var lock = sync.RWMutex{}
	// Grab our router
	router := SetupRouter(db, channelMap, &lock)
	ts := httptest.NewServer(router)
	// defer ts.Close()
	form := url.Values{}
	form.Add("crawl_timeout", "3")
	form.Add("frequency", "3")
	form.Add("failure_threshold", "50")
	form.Add("url", "https://www.khjfcxdoooogle.com/")

	resp, err := http.PostForm(fmt.Sprintf("%s/url", ts.URL), form)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}
	channel, exists := platform.ReadChanFromChanMap("https://www.khjfcxdoooogle.com/", channelMap, &lock)
	if exists {
		channel <- true
	}

	// graceful shutdown

	// channelMap["https://www.khjfcxdoooogle.com/"] <- true
	// Perform a GET request with that handler.
	// w := performRequest(router, "POST", "/url", form)
	// Assert we encoded correctly,
	// the request gives a 200
	// createEntryFunc := handler.CreateURLEntry(db)
	// assert.Equal(t, http.StatusOK, w.Code)
	// // Convert the JSON response to a map
	// var response map[string]string
	// err = json.Unmarshal([]byte(w.Body.String()), &response)
	// // Grab the value & whether or not it exists
	// value, exists := response["url"]
	// // Make some assertions on the correctness of the response.
	// assert.Nil(t, err)
	// assert.True(t, exists)
	// assert.Equal(t, "https://www.khjfcxdoooogle.com/", value)
	ts.Close()
}
