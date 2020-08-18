package dbops

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// DataInDB structure to parse the data in database.
type DataInDB struct {
	UUID             string `json:"uuid"`
	URL              string `json:"url"`
	CrawlTimeout     int    `json:"crawl_timeout"`
	Frequency        int    `json:"frequency"`
	FailureThreshold int    `json:"failure_threshold"`
	IsStatusChecking string `json:"is_status_checking"`
	FailureCount     int    `json:"failure_count"`
}

// // TableExists function to check if table exists
// func TableExists(db *gorm.DB) int {
// 	tableExists, err := db.Query("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'test_database' AND table_name = 'ums'")
// 	if err != nil {
// 		apperrors.ReturnError("Couldn't check if table exists.", err)
// 		fmt.Println("Couldn't check if table exists." + err.Error())
// 		os.Exit(1)
// 	}
// 	var te int
// 	tableExists.Next()
// 	tableExists.Scan(&te)
// 	tableExists.Close()
// 	return te
// }

// // CreateTable function to create a table
// func CreateTable(db *gorm.DB) {
// 	tableCreate, err := db.Query(`CREATE TABLE ums(
// 		uuid varchar(300) NOT NULL,
// 		url varchar(2000) NOT NULL,
// 		crawl_timeout int,
// 		frequency int,
// 		failure_threshold int,
// 		is_status_checking varchar(10),
// 		failure_count int
// 		);`)
// 	if err != nil {
// 		apperrors.ReturnError("Couldn't create a table.", err)
// 		fmt.Println("Couldn't create a table." + err.Error())
// 		// tableCreate.Close()
// 		os.Exit(1)
// 	}
// 	fmt.Println("Table created.")
// 	tableCreate.Close()
// }

// // DropTable to drop the table
// func DropTable(db *gorm.DB) {
// 	dropTable, err := db.Query("DROP TABLE ums;")
// 	if err != nil {
// 		apperrors.ReturnError("Couldn't delete the table.", err)
// 		fmt.Println("Couldn't delete the table." + err.Error())
// 		// dropTable.Close()
// 		os.Exit(1)
// 	}
// 	fmt.Println("Table dropped.")
// 	dropTable.Close()
// }

// FetchURLInfo function to get the data from db if it exists in db
func FetchURLInfo(db *gorm.DB, url string) (DataInDB, bool) {
	var d DataInDB
	dbs := db.Where("url = ?", url).First(&d)
	if dbs.Error != nil {
		fmt.Println(dbs.Error)
		return d, false
	}
	return d, true
}

// InsertURLInfo to insert data in DB
func InsertURLInfo(db *gorm.DB, d DataInDB) error {
	dbi := db.Create(&d)
	// insertQuery := fmt.Sprintf(`INSERT INTO ums VALUES ("%s", "%s", %d, %d, %d, "%s", %d);`, d.UUID, d.URL, d.CrawlTimeout, d.Frequency, d.FailureThreshold, d.IsStatusChecking, d.FailureCount)
	// fmt.Println(insertQuery)
	// insert, err := db.Query(insertQuery)
	// if err != nil {
	// 	fmt.Println("Couldn't insert data in database.", err)
	// 	// insert.Close()
	// 	return err
	// }
	// insert.Close()
	return dbi.Error
}

// UpdateFailureCount function to update failure count
func UpdateFailureCount(db *gorm.DB, url string, failureCount int) error {
	var d DataInDB
	dbufc := db.Model(&d).Where("url = ?", url).Update("failure_count", failureCount)
	// updateQuery := fmt.Sprintf(`UPDATE ums SET %s = %d WHERE url = "%s";`, "failure_count", failureCount, url)
	// err := UpdateURLInfo(db, updateQuery)
	return dbufc.Error
}

// UpdateParams function to update following columns: ["crawl_timeout", "frequency", "failure_threshold"]
func UpdateParams(db *gorm.DB, url string, m map[string]int) error {
	var d DataInDB
	dbup := db.Model(&d).Where("url = ?", url).Updates(m)
	// updateQuery := `UPDATE ums SET`
	// for key, val := range m {
	// 	updateQuery = updateQuery + fmt.Sprintf(` %s = %d`, key, val)
	// }
	// updateQuery = updateQuery + fmt.Sprintf(` failure_count = 0 WHERE url = "%s";`, url)
	// err := UpdateURLInfo(db, updateQuery)
	return dbup.Error
}

// // UpdateServiceStatusUMS to update the status of the URL for UMS service checking
// func UpdateServiceStatusUMS(db *gorm.DB, url string, serviceStatus string) error {
// 	updateQuery := fmt.Sprintf(`UPDATE ums SET is_status_checking = "%s" WHERE url = "%s";`, serviceStatus, url)
// 	err := UpdateURLInfo(db, updateQuery)
// 	return err
// }

// // UpdateURLInfo function to update the information about URL
// func UpdateURLInfo(db *gorm.DB, updateQuery string) error {
// 	update, err := db.Query(updateQuery)
// 	if err != nil {
// 		fmt.Println("Couldn't update data in database.", err)
// 		return err
// 	}
// 	update.Close()
// 	return nil
// }
