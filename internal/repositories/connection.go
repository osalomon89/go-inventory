package repositories

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	DB_HOST = "127.0.0.1"
	DB_PORT = 3306
	DB_NAME = "test-db"
	DB_USER = "root"
	DB_PASS = "secret"
)

var db *sqlx.DB //nolint:gochecknoglobals

func GetConnectionDB() (*sqlx.DB, error) {
	var err error

	if db == nil {
		db, err = sqlx.Connect("mysql", dbConnectionURL())
		if err != nil {
			fmt.Printf("########## DB ERROR: " + err.Error() + " #############")
			return nil, fmt.Errorf("### DB ERROR: %w", err)
		}
	}

	if err := migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func migrate(db *sqlx.DB) error {
	var booksSchema = `
	CREATE TABLE IF NOT EXISTS books (
		id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
		author varchar(200) DEFAULT NULL,
		title longtext,
		price bigint(20) DEFAULT NULL,
		isbn varchar(200) DEFAULT NULL,
		stock bigint(20) DEFAULT NULL,
		created_at datetime(3) DEFAULT NULL,
		updated_at datetime(3) DEFAULT NULL,
		PRIMARY KEY (id),
		UNIQUE KEY isbn (isbn)
	  );`

	_, err := db.Exec(booksSchema)
	if err != nil {
		fmt.Printf("########## DB ERROR: " + err.Error() + " #############")
		return fmt.Errorf("### MIGRATION ERROR: %w", err)
	}

	return nil
}

func dbConnectionURL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True", DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME)
}
