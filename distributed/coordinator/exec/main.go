package main

import (
	"fmt"

	"github.com/snoopy1964/powerPlantMonitor/distributed/coordinator"
)

var dc *coordinator.DatabaseConsumer

func main() {
	ea := coordinator.NewEventAggregator()
	dc = coordinator.NewDatabaseConsumer(ea)
	ql := coordinator.NewQueueListener(ea)

	ch := make(chan string)

	go ql.ListenForNewSources(ch)

	for msg := range ch {
		fmt.Println(msg)
	}

	// var a string
	// fmt.Scanln(&a)
}
