package main

import (
	"flag"
	"fmt"

	"github.com/snoopy1964/powerPlantMonitor/distributed/qutils"
)

var msgBody = flag.String("msgBody", "DuMichAuch", "base64 encoded message body from RabbitMQ-Queues")

func main() {
	flag.Parse()
	fmt.Println(*msgBody)
	fmt.Println(qutils.DecodeMessage(*msgBody))
}
