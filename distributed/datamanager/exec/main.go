package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"

	"github.com/snoopy1964/powerPlantMonitor/distributed/datamanager"
	"github.com/snoopy1964/powerPlantMonitor/distributed/dto"
	"github.com/snoopy1964/powerPlantMonitor/distributed/qutils"
)

// const url = "amqp://guest:guest@localhost:5672"
var url string

func init() {
	ENV_ppm_rabbit_host := os.Getenv("PPM_RABBIT_HOST")
	if ENV_ppm_rabbit_host == "" {
		ENV_ppm_rabbit_host = "localhost"
	}
	ENV_ppm_rabbit_port := os.Getenv("PPM_RABBIT_PORT")
	if ENV_ppm_rabbit_port == "" {
		ENV_ppm_rabbit_port = "5672"
	}
	ENV_ppm_rabbit_user := os.Getenv("PPM_RABBIT_USER")
	if ENV_ppm_rabbit_user == "" {
		ENV_ppm_rabbit_user = "guest"
	}
	ENV_ppm_rabbit_password := os.Getenv("PPM_RABBIT_PASSWORD")
	if ENV_ppm_rabbit_password == "" {
		ENV_ppm_rabbit_password = "guest"
	}

	url = fmt.Sprint(
		"amqp://",
		ENV_ppm_rabbit_user,
		":",
		ENV_ppm_rabbit_password,
		"@",
		ENV_ppm_rabbit_host,
		":",
		ENV_ppm_rabbit_port)
}
func main() {
	conn, ch := qutils.GetChannel(url)
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		qutils.PersistReadingsQueue, // queue string,
		"",                          // consumer string,
		false,                       // autoAck bool,
		true,                        // exclusive bool,
		false,                       // noLocal bool,
		false,                       // noWait bool,
		nil)                         // args amqp.Table)

	if err != nil {
		log.Fatalln("failed to get access to messages")
	}

	for msg := range msgs {
		buf := bytes.NewReader(msg.Body)
		dec := gob.NewDecoder(buf)
		sd := &dto.SensorMessage{}
		dec.Decode(sd)

		err := datamanager.SaveReading(sd)

		if err != nil {
			log.Printf("Failed to save reading from sensor %v. Error: %s", sd.Name, err.Error())
		} else {
			msg.Ack(false)
		}
	}

	// var a string
	// fmt.Scan(&a)

}
