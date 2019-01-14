package main

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/snoopy1964/powerPlantMonitor/distributed/datamanager"
	"github.com/snoopy1964/powerPlantMonitor/distributed/dto"
	"github.com/snoopy1964/powerPlantMonitor/distributed/qutils"
)

const url = "amqp://guest:guest@localhost:5672"

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
