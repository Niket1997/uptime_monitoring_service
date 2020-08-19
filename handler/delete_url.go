package handler

import (
	"fmt"
	"net/http"
	"sync"
	"ums/dbops"
	"ums/platform"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// DeleteURL method to stop the crawling & delete the URL
func DeleteURL(db *gorm.DB, channelMap map[string]chan bool, lock *sync.RWMutex) gin.HandlerFunc {
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
		channel, found := platform.ReadChanFromChanMap(d.URL, channelMap, lock)
		if !found {
			c.JSON(500, gin.H{
				"status": "Error while getting the channel corresponding to URL.",
			})
			return
		}
		channel <- true
		found = platform.DeleteChanFromChannelMap(d.URL, channelMap, lock)
		if found {
			err := dbops.DeleteEntry(db, d.URL)
			if err != nil {
				c.JSON(200, gin.H{
					"status": fmt.Sprintf("Couldn't delete the URL %s.", d.URL),
					"error":  err.Error(),
				})
				return
			}
			c.Status(http.StatusNoContent)
			return
		}
		c.JSON(500, gin.H{
			"status": fmt.Sprintf("Couldn't delete the channel in channelMap."),
		})
	}
}
