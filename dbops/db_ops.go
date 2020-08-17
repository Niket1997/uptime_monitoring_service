package dbops

import (
	"database/sql"
	"fmt"
	"os"
	"ums/handler"
)

// TableExists function to check if table exists
func TableExists(db *sql.DB) int {
	tableExists, err := db.Query("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'test_database' AND table_name = 'ums'")
	if err != nil {
		handler.ReturnError("Couldn't check if table exists.", err)
		fmt.Println("Couldn't check if table exists." + err.Error())
		os.Exit(1)
	}
	var te int
	tableExists.Next()
	tableExists.Scan(&te)
	tableExists.Close()
	return te
}

// CreateTable function to create a table
func CreateTable(db *sql.DB) {
	tableCreate, err := db.Query(`CREATE TABLE ums(
		uuid varchar(30) NOT NULL,
		url varchar(255),
		crawl_timeout int,
		frequency int,
		failure_threshold int,
		is_status_checking varchar(10),
		failure_count int
		);`)
	if err != nil {
		handler.ReturnError("Couldn't create a table.", err)
		fmt.Println("Couldn't create a table." + err.Error())
		os.Exit(1)
	}
	fmt.Println("Table created.")
	tableCreate.Close()
}

// DropTable to drop the table
func DropTable(db *sql.DB) {
	dropTable, err := db.Query("DROP TABLE ums;")
	if err != nil {
		handler.ReturnError("Couldn't delete the table.", err)
		fmt.Println("Couldn't delete the table." + err.Error())
		os.Exit(1)
	}
	fmt.Println("Table dropped.")
	dropTable.Close()
}
