package main

import (
	"sync"
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

	// ChannelMap variable to hold map of channels
	channelMap := make(map[string]chan bool)

	// Lock variable
	var lock = sync.RWMutex{}
	// platform.LockTest(channelMap, &lock)
	// fmt.Println(channelMap)

	r := gin.Default()
	r.POST("/url", handler.GetURLStatus(db, channelMap, &lock))
	r.POST("/urls/:id/deactivate", handler.DeactivateURL(db, channelMap, &lock))
	r.POST("/urls/:id/activate", handler.ActivateURL(db, channelMap, &lock))
	r.DELETE("/urls/:id", handler.DeleteURL(db, channelMap, &lock))
	r.PATCH("/urls/:id", handler.UpdateParams(db, channelMap, &lock))

	r.Run()
}
