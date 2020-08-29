package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"os"
	"sync"
	"ums/dbops"
	"ums/handler"
)

func main() {
	host := "localhost"
	if _, err := os.Stat("/.dockerenv"); err == nil {
		host = "host.docker.internal"
	}

	db, err := gorm.Open("mysql", fmt.Sprintf("root:Anish@6030@tcp(%s:3306)/test_database?charset=utf8&parseTime=True", host))

	if err != nil {
		fmt.Println("database connection failed")
		return
	}
	defer db.Close()

	db.DropTableIfExists(&dbops.DataInDB{})
	db.Debug().AutoMigrate(&dbops.DataInDB{})

	// ChannelMap variable to hold map of channels
	channelMap := make(map[string]chan bool)

	//Lock variable
	var lock = sync.RWMutex{}
	//platform.LockTest(channelMap, &lock)
	fmt.Println(channelMap)
	router := SetupRouter(db, channelMap, &lock)
	err = router.Run()
	if err != nil {
		fmt.Println("error occur")
	}

}

// SetupRouter Method to setup the router
//func SetupRouter() *gin.Engine {
//	r := gin.Default()
//	r.POST("/url", handler.CreateURLEntry())
//	r.GET("/ping", pingEndpoint)
//	return r
//}

func pingEndpoint(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//// SetupRouter method to setup the router
func SetupRouter(db *gorm.DB, channelMap map[string]chan bool, lock *sync.RWMutex) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", pingEndpoint)
	r.POST("/url", handler.CreateURLEntry(db, channelMap, lock))
	r.POST("/urls/:id/deactivate", handler.DeactivateURL(db, channelMap, lock))
	r.POST("/urls/:id/activate", handler.ActivateURL(db, channelMap, lock))
	r.DELETE("/urls/:id", handler.DeleteURL(db, channelMap, lock)) // soft delete
	r.PATCH("/urls/:id", handler.UpdateParams(db, channelMap, lock))
	r.GET("/urls/:id", handler.GetURLInfo(db))
	return r
}
