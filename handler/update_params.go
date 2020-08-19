package handler

import (
	"fmt"
	"sync"
	"time"
	"ums/dbops"
	"ums/platform"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// FormInputURLPatch structure to get input from form
type FormInputURLPatch struct {
	CrawlTimeout     int `form:"crawl_timeout"`
	Frequency        int `form:"frequency"`
	FailureThreshold int `form:"failure_threshold"`
}

// UpdateParams function to update the parameters in DB
func UpdateParams(db *gorm.DB, channelMap map[string]chan bool, lock *sync.RWMutex) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestBody := FormInputURLPatch{}
		err := c.ShouldBind(&requestBody)
		if err != nil {
			c.String(500, err.Error())
		} else {
			id := c.Param("id")
			d, re, err := dbops.FetchURLInfoBasedOnID(db, id)
			if re == 1 {
				c.JSON(500, gin.H{
					"status": "Entry not found in the database.",
					"error":  err.Error(),
				})
				return
			}
			if re == 3 {
				c.JSON(500, gin.H{
					"status": "Error while querying the data in database for existing entry.",
					"error":  err.Error(),
				})
				return
			}
			if d.Status == "Active" {
				channel, found := platform.ReadChanFromChanMap(d.URL, channelMap, lock)
				if !found {
					c.JSON(500, gin.H{
						"status": "Error while getting the channel corresponding to URL.",
					})
					return
				}
				channel <- true
				found = platform.DeleteChanFromChannelMap(d.URL, channelMap, lock)
				fmt.Println("\n\n", channelMap, "\n\n")
				if !found {
					c.JSON(500, gin.H{
						"status": fmt.Sprintf("URL %s was active & couldn't delete the channel.", d.URL),
					})
					return
				}
			}
			m := make(map[string]interface{})
			if requestBody.CrawlTimeout != 0 {
				m["crawl_timeout"] = requestBody.CrawlTimeout
			}
			if requestBody.Frequency != 0 {
				m["frequency"] = requestBody.Frequency
			}
			if requestBody.FailureThreshold != 0 {
				m["failure_threshold"] = requestBody.FailureThreshold
			}
			m["failure_count"] = 0
			m["status"] = "Active"
			err = dbops.UpdateParamsBasedOnURL(db, d.URL, m)
			if err != nil {
				c.JSON(500, gin.H{
					"error": fmt.Sprintf("Couldn't update database for %s. %s", d.URL, err.Error()),
				})
				return
			}
			ticker := time.NewTicker(time.Duration(d.Frequency*1000) * time.Millisecond)
			channel := make(chan bool)
			platform.AddChanToChanMap(d.URL, channelMap, lock, channel)
			go platform.CrawlRoutine(d, db, ticker, channel)
			c.JSON(200, gin.H{
				"status": fmt.Sprintf("The URL %s is updated & started crawling again.", d.URL),
			})
			return

		}
	}
}
