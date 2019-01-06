package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/streadway/amqp"

	"github.com/snoopy1964/powerPlantMonitor/distributed/dto"
	"github.com/snoopy1964/powerPlantMonitor/distributed/qutils"
)

var msgURL = "amqp://guest:guest@localhost:5672"

var name = flag.String("name", "sensor", "name of the sensor")
var freq = flag.Uint("freq", 5, "update frequency in cycles/sec")
var max = flag.Float64("max", 5., "maximum value for generated readings")
var min = flag.Float64("min", 1., "minimum value for generated readings")
var stepsize = flag.Float64("step", 0.1, "maximum allowed change per measurement")

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

var value = r.Float64()*(*max-*min) + *min
var nom = (*max-*min)/2 + *min

func main() {
	flag.Parse()

	conn, ch := qutils.GetChannel(msgURL)
	defer conn.Close()
	defer ch.Close()

	dataQueue := qutils.GetQueue(*name, ch)
	// sensorQueue := qutils.GetQueue(qutils.SensorListQueue, ch)

	msg := amqp.Publishing{Body: []byte(dataQueue.Name)}
	ch.Publish(
		"amq.fanout",
		"",
		false,
		false,
		//		new(amqp.Publishing{Body: []byte(dataQueue.Name),})
		msg)

	dur, _ := time.ParseDuration(strconv.Itoa(1000/int(*freq)) + "ms")

	signal := time.Tick(dur)

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	for range signal {
		calcValue()
		reading := dto.SensorMessage{
			Name:  *name,
			Value: value,
			Tst:   time.Now(),
		}

		buf.Reset()
		enc.Encode(reading)

		msg := amqp.Publishing{
			Body: buf.Bytes(),
		}

		ch.Publish("", dataQueue.Name, false, false, msg)

		log.Printf("Reading sent. Value: %v\n", value)
	}
}

func calcValue() {
	var maxStep, minStep float64
	if value < nom {
		maxStep = *stepsize
		minStep = -1 * *stepsize * (value - *min) / (nom - *min)
	} else {
		minStep = -1 * *stepsize
		maxStep = *stepsize * (*max - value) / (*max - nom)
	}

	value += r.Float64()*(maxStep-minStep) + minStep
}