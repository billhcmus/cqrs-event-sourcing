package eventsourcing

import (
	"github.com/jinzhu/gorm"
)

// DB to access to DB in app
var DB *gorm.DB

// Tx is a alias of *gorm.DB
type Tx = *gorm.DB

// Init initialize the db
func Init(dbURL string) error {
	db, err := gorm.Open("postgres", dbURL)

	if err != nil {
		return err
	}

	DB = db
	return DB.DB().Ping()
}