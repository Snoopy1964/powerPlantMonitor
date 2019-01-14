package datamanager

import (
	"errors"
	"log"

	"github.com/snoopy1964/powerPlantMonitor/distributed/dto"
)

var sensors map[string]int

func SaveReading(reading *dto.SensorMessage) error {
	if sensors[reading.Name] == 0 {
		getSensors()
	}

	if sensors[reading.Name] == 0 {
		log.Printf("------ Message: %v", reading)
		return errors.New("Unable to find sensor for name '" + reading.Name)
	}

	s := `
		INSERT INTO sensor_reading
			(value, sensor_id, taken_on)
		VALUES
			($1, $2, $3)
	`
	_, err := db.Exec(s, reading.Value, sensors[reading.Name], reading.Tst)

	return err
}

func getSensors() {
	sensors = make(map[string]int)
	q := `
		SELECT id, 
		name FROM sensor
	`

	rows, _ := db.Query(q)

	for rows.Next() {
		var id int
		var name string

		rows.Scan(&id, &name)

		sensors[name] = id
	}
}
