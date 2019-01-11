package datamanager

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://user:password@localhost/powerPlantMonitor?sslmode=disable")
	if err != nil {
		log.Fatalln(fmt.Errorf("Unable to connect to database: %v", err))
	}
	p := db.Ping()
	log.Printf("ping from DB: %v", p)

}
