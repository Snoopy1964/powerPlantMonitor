package qutils

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/snoopy1964/powerPlantMonitor/distributed/dto"

	"github.com/streadway/amqp"
)

// queue for discovering sensors
const SensorDiscoveryExchange = "SensorDiscovery"

// queue for persistent readings
const PersistReadingsQueue = "PersistReading"

func GetChannel(url string) (*amqp.Connection, *amqp.Channel) {
	log.Println("connection to Messagebroker with connection string: ", url)
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to establish a connection to message broker")
	ch, err := conn.Channel()
	failOnError(err, "Failed to get channel for connection")

	return conn, ch
}

func GetQueue(name string, ch *amqp.Channel, autoDelete bool) *amqp.Queue {
	q, err := ch.QueueDeclare(name, false, false, autoDelete, false, nil)
	failOnError(err, "Failed to declare a queue")

	return &q
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func DecodeMessage(msgBody64 string) dto.SensorMessage {
	m, err := base64.StdEncoding.DecodeString(msgBody64)
	if err != nil {
		log.Fatalln("Base64 decoding of message failed: ", err)
	}
	r := bytes.NewReader(m)
	d := gob.NewDecoder(r)

	sd := new(dto.SensorMessage)
	d.Decode(sd)

	return *sd
}
