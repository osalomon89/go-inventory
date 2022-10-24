package repositories

import (
	"fmt"

	"github.com/osalomon89/go-inventory/internal/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	DB_HOST = "127.0.0.1"
	DB_PORT = 3306
	DB_NAME = "test-db"
	DB_USER = "root"
	DB_PASS = "secret"
)

// gorm
func GetConnectionDB() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dbConnectionURL()))
	if err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	if err := migrate(db); err != nil {
		return nil, err
	}
	return db, nil
}
func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&domain.Book{})
	if err != nil {
		return fmt.Errorf("error migrating db: %w", err)
	}
	return nil
}
func dbConnectionURL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True", DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME)
}

/*
var db *sqlx.DB

func GetConnectionDB() (*sqlx.DB, error) {
	var err error
	if db == nil {
		db, err = sqlx.Connect("mysql", dbConnectionURL())
		if err != nil {
			fmt.Printf("#####DB ERROR: " + err.Error() + "#####")
			return nil, fmt.Errorf("### db error: %w", err)
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
		fmt.Printf("#####DB ERROR: " + err.Error() + "#####")
		return fmt.Errorf("### db error: %w", err)
	}
	return nil
}
*/
