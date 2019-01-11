package coordinator

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

	"github.com/Snoopy1964/powerPlantMonitor/distributed/dto"
	"github.com/snoopy1964/powerPlantMonitor/distributed/qutils"
	"github.com/streadway/amqp"
)

// maximum rate for throtteling message rate, -> only every 5 seconds
const maxRate = 5 * time.Second

type DatabaseConsumer struct {
	er      EventRaiser
	conn    *amqp.Connection
	ch      *amqp.Channel
	queue   *amqp.Queue
	sources []string
}

func NewDatabaseConsumer(er EventRaiser) *DatabaseConsumer {
	dc := DatabaseConsumer{
		er: er,
	}

	dc.conn, dc.ch = qutils.GetChannel(url)
	dc.queue = qutils.GetQueue(qutils.PersistReadingsQueue, dc.ch, false)

	dc.er.AddListener("DataSourceDiscovered", func(eventData interface{}) {
		dc.SubscribeToDataEvent(eventData.(string))
	})

	return &dc
}

func (dc *DatabaseConsumer) SubscribeToDataEvent(eventName string) {
	for _, v := range dc.sources {
		if v == eventName {
			return
		}
	}

	dc.er.AddListener("MessageReceived_"+eventName, func() func(interface{}) {
		prevTime := time.Unix(0, 0)

		buf := new(bytes.Buffer)

		return func(eventData interface{}) {
			ed := eventData.(EventData)
			// log.Printf("!!!! Event ed !!!!: %v", ed)
			if time.Since(prevTime) > maxRate {
				prevTime = time.Now()

				sm := dto.SensorMessage{
					Name:  ed.Name,
					Value: ed.Value,
					Tst:   ed.Tst,
				}

				log.Printf("!!!! Event sm!!!!: %v", sm)

				buf.Reset()

				enc := gob.NewEncoder(buf)
				enc.Encode(sm)

				msg := amqp.Publishing{
					Body: buf.Bytes(),
				}

				dc.ch.Publish(
					"",                          // exchange string,
					qutils.PersistReadingsQueue, // key string,
					false,                       // mandatory bool,
					false,                       // immediate bool,
					msg)                         //msg amqp.Publishing)
			}
		}
	}())
}
