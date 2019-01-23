package qutils

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	"github.com/snoopy1964/powerPlantMonitor/distributed/dto"

	"github.com/streadway/amqp"
)

// queue for discovering sensors
const SensorDiscoveryExchange = "SensorDiscovery"

// queue for persistent readings
const PersistReadingsQueue = "PersistReading"

func connectToRabbitMQ(url string) (conn *amqp.Connection, err error) {
	t0 := time.Now()
	var err0 error
	for {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		} else {
			if err0 == nil {
				log.Printf("%s\n", err.Error())
			} else {
				if err0.Error() != err.Error() {
					log.Printf("%s\n", err.Error())
				}
			}
			err0 = err
		}

		time.Sleep(1000 * time.Millisecond)
		log.Printf("Trying to reconnect to RabbitMQ at %s\n", url)
		if time.Now().Sub(t0) > 30*time.Second {
			log.Println("Timeout reconnecting....")
			break
		}
		err0 = err
	}
	return conn, err
}

func GetChannel(url string) (*amqp.Connection, *amqp.Channel) {
	log.Println("connection to Messagebroker with connection string: ", url)
	conn, err := connectToRabbitMQ(url)
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
