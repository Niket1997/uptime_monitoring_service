package main

import (
	"database/sql"
	"ums/apperrors"
	"ums/dbops"
	"ums/handler"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// u1 := uuid.Must(uuid.NewV4())
	// fmt.Printf("UUIDv4: %s\n", u1)

	// // or error handling
	// u2, err := uuid.NewV4()
	// if err != nil {
	// 	fmt.Printf("Something went wrong: %s", err)
	// 	return
	// }
	// fmt.Printf("UUIDv4: %s\n", u2)

	// // // Parsing UUID from string input
	// // u2, err := uuid.FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	// // if err != nil {
	// // 	fmt.Printf("Something went wrong: %s", err)
	// // 	return
	// // }
	// // fmt.Printf("Successfully parsed: %s", u2)

	db, err := sql.Open("mysql", "root:Anish@6030@tcp(127.0.0.1:3306)/test_database")

	if err != nil {
		apperrors.ReturnError("Couldn't create database connection.", err)
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
