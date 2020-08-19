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

// ActivateURL function to activate the inactive URL
func ActivateURL(db *gorm.DB, channelMap map[string]chan bool, lock *sync.RWMutex) gin.HandlerFunc {
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
		if d.Status == "Active" {
			c.JSON(500, gin.H{
				"status": "The URL is already active.",
			})
			return
		}
		// channel, found := platform.ReadChanFromChanMap(d.URL, channelMap, lock)
		// if !found {
		// 	c.JSON(500, gin.H{
		// 		"status": "Error while getting the channel corresponding to URL.",
		// 	})
		// 	return
		// }
		ticker := time.NewTicker(time.Duration(d.Frequency*1000) * time.Millisecond)
		channel := make(chan bool)
		platform.AddChanToChanMap(d.URL, channelMap, lock, channel)
		go platform.CrawlRoutine(d, db, ticker, channel)
		c.JSON(200, gin.H{
			"status": fmt.Sprintf("The URL %s is activated.", d.URL),
		})
	}
}
