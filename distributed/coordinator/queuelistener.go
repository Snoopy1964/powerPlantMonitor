package coordinator

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"

	"github.com/snoopy1964/powerPlantMonitor/distributed/dto"
	"github.com/snoopy1964/powerPlantMonitor/distributed/qutils"
	"github.com/streadway/amqp"
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

type QueueListener struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	sources map[string]<-chan amqp.Delivery
	ea      *EventAggregator
}

func NewQueueListener(ea *EventAggregator) *QueueListener {
	ql := QueueListener{
		sources: make(map[string]<-chan amqp.Delivery),
		ea:      ea,
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

	ql.ch.Publish(qutils.SensorDiscoveryExchange, "", false, false, amqp.Publishing{})
}

func (ql *QueueListener) ListenForNewSources() {
	q := qutils.GetQueue("", ql.ch, true)
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
			fmt.Printf("new source discovered: %s\n", string(msg.Body))
			ql.ea.PublishEvent("DataSourceDiscovered", string(msg.Body))
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

		fmt.Printf("Message received: %v\n", sd)

		ed := EventData{
			Name:  sd.Name,
			Tst:   sd.Tst,
			Value: sd.Value,
		}
		ql.ea.PublishEvent("MessageReceived_"+msg.RoutingKey, ed)
	}
}
