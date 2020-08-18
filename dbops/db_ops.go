package dbops

import (
	"database/sql"
	"fmt"
	"os"
	"ums/apperrors"

	"github.com/gofrs/uuid"
)

// DataInDB structure to parse the data in database.
type DataInDB struct {
	UUID             uuid.UUID `json:"uuid"`
	URL              string    `json:"url"`
	CrawlTimeout     int       `json:"crawl_timeout"`
	Frequency        int       `json:"frequency"`
	FailureThreshold int       `json:"failure_threshold"`
	IsStatusChecking string    `json:"is_status_checking"`
	FailureCount     int       `json:"failure_count"`
}

// TableExists function to check if table exists
func TableExists(db *sql.DB) int {
	tableExists, err := db.Query("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'test_database' AND table_name = 'ums'")
	if err != nil {
		apperrors.ReturnError("Couldn't check if table exists.", err)
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
		uuid varchar(300) NOT NULL,
		url varchar(2000) NOT NULL,
		crawl_timeout int,
		frequency int,
		failure_threshold int,
		is_status_checking varchar(10),
		failure_count int
		);`)
	if err != nil {
		apperrors.ReturnError("Couldn't create a table.", err)
		fmt.Println("Couldn't create a table." + err.Error())
		// tableCreate.Close()
		os.Exit(1)
	}
	fmt.Println("Table created.")
	tableCreate.Close()
}

// DropTable to drop the table
func DropTable(db *sql.DB) {
	dropTable, err := db.Query("DROP TABLE ums;")
	if err != nil {
		apperrors.ReturnError("Couldn't delete the table.", err)
		fmt.Println("Couldn't delete the table." + err.Error())
		// dropTable.Close()
		os.Exit(1)
	}
	fmt.Println("Table dropped.")
	dropTable.Close()
}

// FetchURLInfo function to get the data from db if it exists in db
func FetchURLInfo(db *sql.DB, url string) (DataInDB, bool) {
	var d DataInDB
	err := db.QueryRow("SELECT * FROM ums where url = ?", url).Scan(&d.UUID, &d.URL, &d.CrawlTimeout, &d.Frequency, &d.FailureThreshold, &d.IsStatusChecking, &d.FailureCount)
	if err != nil {
		fmt.Println(err)
		return d, false
	}
	return d, true
}

// InsertURLInfo to insert data in DB
func InsertURLInfo(db *sql.DB, d DataInDB) error {
	insertQuery := fmt.Sprintf(`INSERT INTO ums VALUES ("%s", "%s", %d, %d, %d, "%s", %d);`, d.UUID, d.URL, d.CrawlTimeout, d.Frequency, d.FailureThreshold, d.IsStatusChecking, d.FailureCount)
	fmt.Println(insertQuery)
	insert, err := db.Query(insertQuery)
	if err != nil {
		fmt.Println("Couldn't insert data in database.", err)
		// insert.Close()
		return err
	}
	insert.Close()
	return nil
}

// UpdateFailureCount function to update failure count
func UpdateFailureCount(db *sql.DB, url string, failureCount int) error {
	updateQuery := fmt.Sprintf(`UPDATE ums SET %s = %d WHERE url = "%s";`, "failure_count", failureCount, url)
	err := UpdateURLInfo(db, updateQuery)
	return err
}

// UpdateParams function to update following columns: ["crawl_timeout", "frequency", "failure_threshold"]
func UpdateParams(db *sql.DB, url string, m map[string]int) error {
	updateQuery := `UPDATE ums SET`
	for key, val := range m {
		updateQuery = updateQuery + fmt.Sprintf(` %s = %d`, key, val)
	}
	updateQuery = updateQuery + fmt.Sprintf(` failure_count = 0 WHERE url = "%s";`, url)
	err := UpdateURLInfo(db, updateQuery)
	return err
}

// UpdateServiceStatusUMS to update the status of the URL for UMS service checking
func UpdateServiceStatusUMS(db *sql.DB, url string, serviceStatus string) error {
	updateQuery := fmt.Sprintf(`UPDATE ums SET is_status_checking = "%s" WHERE url = "%s";`, serviceStatus, url)
	err := UpdateURLInfo(db, updateQuery)
	return err
}

// UpdateURLInfo function to update the information about URL
func UpdateURLInfo(db *sql.DB, updateQuery string) error {
	update, err := db.Query(updateQuery)
	if err != nil {
		fmt.Println("Couldn't update data in database.", err)
		return err
	}
	update.Close()
	return nil
}
