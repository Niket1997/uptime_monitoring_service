package handler

import (
	"database/sql"
	"fmt"
	"net/http"
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

func checkIfURLExistsInDb(db *sql.DB) {

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
				ReturnError("Error while generating UUID.", err)
				return
			}
			insert, err := db.Query("INSERT INTO ums VALUES ( 5, 'Mhataru' )")
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"status": "Error while inserting data into table.",
					"error":  err.Error(),
				})
			}
			defer insert.Close()

			status := platform.GetRequest(url, timeout)
			c.JSON(200, gin.H{
				"id":                id,
				"url":               url,
				"crawl_timeout":     timeout,
				"frequency":         frequency,
				"failure_threshold": failureThreshold,
				"status":            status,
				"failure_count":     1,
			})

		}
	}
}
