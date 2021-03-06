package dto

import (
	"encoding/gob"
	"time"
)

type SensorMessage struct {
	Name  string
	Value float64
	Tst   time.Time
}

func init() {
	gob.Register(SensorMessage{})
}
