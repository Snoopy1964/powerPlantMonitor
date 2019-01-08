package main

import (
	"fmt"

	"github.com/Snoopy1964/powerPlantMonitor/distributed/coordinator"
)

func main() {
	ql := coordinator.NewQueueListener()
	go ql.ListenForNewSources()

	var a string
	fmt.Scanln(&a)
}
