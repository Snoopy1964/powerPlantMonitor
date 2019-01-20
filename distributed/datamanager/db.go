package datamanager

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB
var dbConUrl string

func init() {
	// get connection string from environment variable
	dbConUrl := os.Getenv("PPM_DB_CONNECT")
	if dbConUrl == "" {
		dbConUrl = "postgres://user:password@localhost/powerPlantMonitor?sslmode=disable"
	}
	log.Println("connection to Messagebroker with connection string: ", url)

	var err error
	db, err = sql.Open("postgres", dbConUrl)
	if err != nil {
		log.Fatalln(fmt.Errorf("Unable to connect to database: %v", err))
	}
	p := db.Ping()
	log.Printf("ping from DB: %v", p)

}
