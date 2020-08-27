package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"sync"
	"time"
	"ums/dbops"
	"ums/platform"
)

// FormInputURLStatus structure to get input from form
type FormInputURLStatus struct {
	URL              string `form:"url" binding:"required"`
	CrawlTimeout     int    `form:"crawl_timeout" binding:"required"`
	Frequency        int    `form:"frequency" binding:"required"`
	FailureThreshold int    `form:"failure_threshold" binding:"required"`
}

// CreateURLEntry function to get the status of the URL
func CreateURLEntry(db *gorm.DB, channelMap map[string]chan bool, lock *sync.RWMutex) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestBody := FormInputURLStatus{}
		err := c.ShouldBind(&requestBody)
		if err != nil {
			c.String(500, err.Error())
		} else {
			url := requestBody.URL
			timeout := requestBody.CrawlTimeout
			frequency := requestBody.Frequency
			failureThreshold := requestBody.FailureThreshold
			// add validations

			d, re, err := dbops.FetchURLInfoBasedOnURL(db, url)
			// there is no need of re
			if re == 3 {
				c.JSON(500, gin.H{
					"status": "Error while querying the data in database for existing entry.",
					"error":  err.Error(),
				})
				return
			}
			createGoRoutine := false
			if re == 1 {
				// gorm callbacks (won't have to write)
				id, err := uuid.NewV1()
				if err != nil {
					c.JSON(500, gin.H{
						"status": "Error while generating UUID.",
						"error":  err.Error(),
					})
					return
				}
				d = dbops.DataInDB{
					UUID:             fmt.Sprintf("%s", id),
					URL:              url,
					CrawlTimeout:     timeout,
					Frequency:        frequency,
					FailureThreshold: failureThreshold,
					Status:           "Active", // enum
					FailureCount:     0,
				}
				err = dbops.InsertURLInfo(db, d)
				if err != nil {
					c.JSON(500, gin.H{
						"status": "Error while inserting data in database.",
						"error":  err.Error(),
					})
					return
				}
				createGoRoutine = true
			}

			c.JSON(200, gin.H{
				"id":                d.UUID,
				"url":               d.URL,
				"crawl_timeout":     d.CrawlTimeout,
				"frequency":         d.Frequency,
				"failure_threshold": d.FailureThreshold,
				"status":            d.Status,
				"failure_count":     d.FailureCount,
			})

			if createGoRoutine {
				ticker := time.NewTicker(time.Duration(d.Frequency*1000) * time.Millisecond)
				channel := make(chan bool)
				platform.AddChanToChanMap(d.URL, channelMap, lock, channel)
				go platform.CrawlRoutine(d, db, ticker, channel)
			}
		}
	}
}
