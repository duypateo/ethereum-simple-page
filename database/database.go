package database

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	driver   = "mysql"
	username = "root"
	password = "root"
	port     = "3306"
	host     = "127.0.0.1"
	schema   = "ethereum"
	// DB - for db instance
	DB *sql.DB
)

// History represent for table history in database
type History struct {
	Address  string `db:"address"`
	Datetime string `db:"datetime"`
}

// InitConnection - connect db
func InitConnection() error {
	conString := os.Getenv("CLEARDB_DATABASE_URL")
	if conString == "" {
		conString := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + schema
	}

	var err error
	DB, err = sql.Open(driver, conString)
	if err != nil {
		panic(err)
	}

	return nil
}

// InsertToHistory - insert to db
func InsertToHistory(addr string) bool {
	if IsExistAddress(addr) {
		return UpdateToHistory(addr)
	}

	insertQuery := "INSERT INTO histories VALUES ('" + addr + "', now())"
	insert, err := DB.Query(insertQuery)
	if err != nil {
		panic(err)
	}

	defer insert.Close()
	return true
}

// UpdateToHistory - update to db
func UpdateToHistory(addr string) bool {
	updateQuery := "UPDATE histories SET datetime = now() WHERE address = '" + addr + "'"
	update, err := DB.Query(updateQuery)
	if err != nil {
		panic(err)
	}

	defer update.Close()
	return true
}

// GetHistories - get history record from db
func GetHistories() []History {
	var histories []History

	selectQuery := `SELECT address, datetime FROM histories ORDER BY datetime DESC`
	rows, err := DB.Query(selectQuery)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		history := History{}
		err = rows.Scan(&history.Address, &history.Datetime)
		if err != nil {
			panic(err)
		}
		histories = append(histories, history)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return histories
}

// IsExistAddress - check if address is stored
func IsExistAddress(addr string) bool {
	var address string

	sqlStatement := "SELECT address FROM histories WHERE address='" + addr + "'"
	row := DB.QueryRow(sqlStatement)

	err := row.Scan(&address)
	if err != nil {
		return false
	}

	return true
}

// IsExist - insert to db
func IsExist(addr string) error {
	var address string

	sqlStatement := "SELECT address FROM histories WHERE address='" + addr + "'"
	row := DB.QueryRow(sqlStatement)

	err := row.Scan(&address)
	if err != nil {
		return err
	}

	return nil
}
