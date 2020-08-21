package handler

import (
	"ums/dbops"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetURLInfo function to get the information about URL
func GetURLInfo(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		c.JSON(200, gin.H{
			"id":                d.UUID,
			"url":               d.URL,
			"crawl_timeout":     d.CrawlTimeout,
			"frequency":         d.Frequency,
			"failure_threshold": d.FailureThreshold,
			"status":            d.Status,
			"failure_count":     d.FailureCount,
		})
	}
}
