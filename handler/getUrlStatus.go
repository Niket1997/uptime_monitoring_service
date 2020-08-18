package handler

import (
	"database/sql"
	"net/http"
	"ums/dbops"
	"ums/platform"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// FormInputURLStatus structure to get input from form
type FormInputURLStatus struct {
	URL              string `form:"url" binding:"required"`
	CrawlTimeout     int    `form:"crawl_timeout" binding:"required"`
	Frequency        int    `form:"frequency" binding:"required"`
	FailureThreshold int    `form:"failure_threshold" binding:"required"`
}

// GetURLStatus function to get the status of the URL
func GetURLStatus(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestBody := FormInputURLStatus{}
		err := c.ShouldBind(&requestBody)
		if err != nil {
			c.String(http.StatusOK, err.Error())
		} else {
			url := requestBody.URL
			timeout := requestBody.CrawlTimeout
			frequency := requestBody.Frequency
			failureThreshold := requestBody.FailureThreshold
			id, err := uuid.NewV1()
			if err != nil {
				c.JSON(500, gin.H{
					"status": "Error while generating UUID.",
					"error":  err.Error(),
				})
				// apperrors.ReturnError("Error while generating UUID.", err)
				return
			}
			status := platform.GetRequest(url, timeout)

			d, re := dbops.FetchURLInfo(db, url)
			failureCount := 0
			if status == "Inactive" && re == true {
				failureCount = d.FailureCount + 1
				err := dbops.UpdateFailureCount(db, url, failureCount)
				if err != nil {
					c.JSON(500, gin.H{
						"status": "Error while updating data in database.",
						"error":  err.Error(),
					})
				}
			} else {

				if re == false {
					if status == "Inactive" {
						failureCount = 1
					}
					dataDB := dbops.DataInDB{
						UUID:             id,
						URL:              url,
						CrawlTimeout:     timeout,
						Frequency:        frequency,
						FailureThreshold: failureThreshold,
						IsStatusChecking: "Active",
						FailureCount:     failureCount,
					}
					err := dbops.InsertURLInfo(db, dataDB)
					if err != nil {
						c.JSON(500, gin.H{
							"status": "Error while inserting data in database.",
							"error":  err.Error(),
						})
					}
				}
			}

			c.JSON(200, gin.H{
				"id":                id,
				"url":               url,
				"crawl_timeout":     timeout,
				"frequency":         frequency,
				"failure_threshold": failureThreshold,
				"status":            status,
				"failure_count":     failureCount,
			})

		}
	}
}
