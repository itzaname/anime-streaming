package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

var DB *sql.DB

// Connect establishes a connection to a mysql database and sets DB to the instance.
// Input is in DNS format
// username:password@protocol(address)/dbname?param=value
func Connect(dsn string) error {
	db, err := sql.Open("mysql",dsn)
	if err != nil {
		return err
	}

	DB = db
	return db.Ping()
}