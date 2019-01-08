package coordinator

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/Snoopy1964/powerPlantMonitor/distributed/dto"

	"github.com/snoopy1964/powerPlantMonitor/distributed/qutils"
	"github.com/streadway/amqp"
)

const url = "amqp://guest:guest@localhost:5672"

type QueueListener struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	sources map[string]<-chan amqp.Delivery
	ea      *EventAggregator
}

func NewQueueListener() *QueueListener {
	ql := QueueListener{
		sources: make(map[string]<-chan amqp.Delivery),
		ea:      NewEventAggregator(),
	}

	ql.conn, ql.ch = qutils.GetChannel(url)

	return &ql
}

func (ql *QueueListener) DiscoverSensors() {
	ql.ch.ExchangeDeclare(
		qutils.SensorDiscoveryExchange, // name string,
		"fanout",                       // kind string,
		false,                          // durable bool,
		false,                          // autoDelete bool,
		false,                          // internal bool,
		false,                          // noWait bool,
		nil)                            // args amqp.Table)
}

func (ql *QueueListener) ListenForNewSources() {
	q := qutils.GetQueue("", ql.ch)
	ql.ch.QueueBind(
		q.Name,
		"",
		"amq.fanout",
		false,
		nil)

	msgs, _ := ql.ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)

	ql.DiscoverSensors()

	fmt.Println("Listening for new sensors")

	for msg := range msgs {
		if ql.sources[string(msg.Body)] == nil {
			sourceChan, _ := ql.ch.Consume(
				string(msg.Body), //queue string,
				"",               //consumer string,
				true,             //autoAck bool,
				false,            //exclusive bool,
				false,            //noLocal bool,
				false,            //noWait bool,
				nil)              //args amqp.Table)

			ql.sources[string(msg.Body)] = sourceChan
			
			go ql.AddListener(sourceChan)
			fmt.Printf("new ListenerQueue added: %s\n:", string(msg.Body))
		}
	}

}

func (ql *QueueListener) AddListener(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		r := bytes.NewReader(msg.Body)
		d := gob.NewDecoder(r)
		sd := new(dto.SensorMessage)
		d.Decode(sd)
		// fmt.Printf("Message received (undecoded): %v\n", msg)
		fmt.Printf("Message received: %v\n", sd)

		ed := EventData{
			Name:  sd.Name,
			Tst:   sd.Tst,
			Value: sd.Value,
		}
		ql.ea.PublishEvent("MessageReceived_"+msg.RoutingKey, ed)
	}
}
