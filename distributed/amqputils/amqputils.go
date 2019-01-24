package amqputils

import (
	"context"
	"log"

	"github.com/streadway/amqp"
)

type message []byte

// Struct embedding
type session struct {
	*amqp.Connection
	*amqp.Channel
}

// Close the session if the connection is closed
func (s session) Close() error {
	if s.Connection == nil {
		return nil
	}

	return s.Connection.Close()
}

//------------------------------- https://www.goin5minutes.com/blog/channel_over_channel/
func Redial(ctx context.Context, url string) chan chan session {
	sessions := make(chan chan session)

	go func() {
		s := make(chan session)
		defer close(sessions)

		for {
			select {
			case sessions <- s:
			case <-ctx.Done():
				log.Println("shutting down session factory...")
				return
			}
			conn, err := amqp.Dial(url)
			if err != nil {
				log.Fatalf("cannot (re)dial: %v; %q", err, url)
			}

			ch, err := conn.Channel()
			if err != nil {
				log.Fatalf("cannot create channel: %v", err)
			}

			/* do I need that here?
			err := ch.ExchangeDeclare(
				exchange, // name string,
				"fanout", //kind string,
				false,    // durable bool,
				true,     // autoDelete bool,
				false,    // internal bool,
				false,    // noWait bool,
				nil)      // args amqp.Table)
			*/

			select {
			case s <- session{conn, ch}:
			case <-ctx.Done():
				log.Println("closing new session")
			}
		}

	}()

	return sessions
}
