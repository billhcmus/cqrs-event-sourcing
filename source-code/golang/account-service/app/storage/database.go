package main

import (
	"log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	//db, err := gorm.Open("postgres", "host=127.0.0.1 port=5432 user=postgres dbname=eventstore password=1234 sslmode=disable")
	db, err := gorm.Open("postgres", "postgres://postgres:1234@localhost/eventstore?sslmode=disable")
	
	if err != nil {
		panic(err)
	} else {
		log.Println("Connect to database successfully")
	}
	defer db.Close()
	
	err = db.Exec(`
	 	CREATE TABLE users (
			id UUID NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL,
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
			deleted_at TIMESTAMP WITH TIME ZONE NOT NULL,
			version BIGINT NOT NULL,

			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			balance BIGINT NOT NULL,

			PRIMARY KEY (id)
		 );
	  `).Error

	if err != nil {
		panic(err)
	} else {
		log.Println("Created table USERS")
	}

	err = db.Exec(`
		CREATE TABLE users_events (
			id UUID NOT NULL,
			timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
			aggregate_id UUID NOT NULL,
			aggregate_type TEXT NOT NULL,
			action TEXT NOT NULL,
			version BIGINT NOT NULL,
			type TEXT NOT NULL,
			data JSONB NOT NULL,
			
			PRIMARY KEY (id)
		)
	`).Error

	if err != nil {
		panic(err)
	} else {
		log.Println("Created table USERS_EVENTS")
	}
}