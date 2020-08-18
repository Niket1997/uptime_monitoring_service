package main

import (
	"ums/apperrors"
	"ums/dbops"
	"ums/handler"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func main() {

	db, err := gorm.Open("mysql", "root:Anish@6030@tcp(127.0.0.1:3306)/test_database?charset=utf8&parseTime=True")

	if err != nil {
		apperrors.ReturnError("Couldn't create database connection.", err)
		return
	}
	defer db.Close()

	db.DropTableIfExists(&dbops.DataInDB{})
	db.Debug().AutoMigrate(&dbops.DataInDB{})
	// id, _ := uuid.NewV1()
	// entry := dbops.DataInDB{
	// 	UUID:             fmt.Sprintf("%s", id),
	// 	URL:              "https://google.com/",
	// 	CrawlTimeout:     10,
	// 	Frequency:        30,
	// 	FailureThreshold: 10,
	// 	IsStatusChecking: "Active",
	// 	FailureCount:     1,
	// }

	// db.Create(&entry)

	// te := dbops.TableExists(db)
	// if te == 0 {
	// 	dbops.CreateTable(db)
	// } else {
	// 	dbops.DropTable(db)
	// 	dbops.CreateTable(db)
	// }

	r := gin.Default()
	r.GET("/url", handler.GetURLStatus(db))
	r.Run()
}
