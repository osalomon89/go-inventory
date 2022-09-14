package repositories

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/osalomon89/go-inventory/internal/domain"
)

const (
	DB_HOST = "127.0.0.1"
	DB_PORT = 3306
	DB_NAME = "test-db"
	DB_USER = "root"
	DB_PASS = "secret"
)

var db *gorm.DB
var booksSchema *domain.Book

func GetConnectionDB() (*gorm.DB, error) {
	var err error

	if db == nil {
		db, err = gorm.Open("mysql", dbConnectionURL())
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

func migrate(db *gorm.DB) error {
	var err = db.AutoMigrate(booksSchema)
	/*
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
	*/
	if err.Error != nil {
		fmt.Printf("########## DB ERROR #############")
		return fmt.Errorf("### MIGRATION ERROR:")
	}

	return nil
}

func dbConnectionURL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True", DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME)
}
