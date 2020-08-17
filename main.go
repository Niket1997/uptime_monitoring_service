package main

import (
	"database/sql"
	"ums/handler"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:Anish@6030@tcp(127.0.0.1:3306)/test_database")

	if err != nil {
		handler.ReturnError("Couldn't create database connection.", err)
		return
	}
	defer db.Close()

	te := dbops.TableExists(db)
	if te == 0 {
		dbops.CreateTable(db)
	} else {
		dbops.DropTable(db)
		dbops.CreateTable(db)
	}

	r := gin.Default()
	r.GET("/url", handler.GetURLStatus(db))
	r.Run()
}
