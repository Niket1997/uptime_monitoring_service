package platform

import (
	"fmt"
	"time"
	"ums/dbops"

	"github.com/jinzhu/gorm"
)

// CrawlRoutine function to crawl the create the go routine for a URL
func CrawlRoutine(d dbops.DataInDB, db *gorm.DB, ticker *time.Ticker, channel chan bool) {
	failureCount := 0
	for {
		select {
		case <-channel:
			{
				fmt.Printf("Stopped crawling %s as per your request.\n", d.URL)
				err := dbops.UpdateStatus(db, d.URL, "Inactive")
				if err != nil {
					fmt.Println("Failed to set the database status of URL to Inactive. ", err.Error())
				}
				return
			}
		case <-ticker.C:
			{
				status := GetRequest(d.URL, d.CrawlTimeout)
				if status == "Inactive" {
					failureCount++
					err := dbops.UpdateFailureCount(db, d.URL)
					if err != nil {
						fmt.Println("Stopped crawling as update to failure count failed. ", err.Error())
						return
					}
					if failureCount == d.FailureThreshold {
						fmt.Printf("Stopped crawling %s as failure count exceeded frequency threshold.\n", d.URL)
						err = dbops.UpdateStatus(db, d.URL, "Inactive")
						if err != nil {
							fmt.Println("Failed to set the database status of URL to Inactive. ", err.Error())
						}
						return
					}
				}
			}
		}
	}

}
