package coordinator

import (
	"github.com/snoopy1964/powerPlantMonitor/distributed/qutils"
	"github.com/streadway/amqp"
)

const url = "ampq://guest:guest@localhost:5672"

type QueueListener struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewQueueListener() * QueueListener {
	ql := QueueListener {}

	ql.conn, ql.ch := qutils.GetChannel(url)

	return &ql
}

func (ql *QueueListener) ListenForNewSources() {
	q:= qutils.GetQueue("", ch)
	ql.ch.QueueBind	(
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
			flase,
			nil)
	
			
}