package main

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/Snoopy1964/powerPlantMonitor/distributed/dto"
)

func main() {
	msgBody := `N/+BAwEBDVNlbnNvck1lc3NhZ2UB/4IAAQMBBE5hbWUBDAABBVZhbHVlAQgAAQNUc3QB/4QAAAAQ/4MFAQEEVGltZQH/hAAAADP/ggETYm9pbGVyX3ByZXNz
dXJlX291dAH4y/ZJPQi8CUABDwEAAAAO08qKAzVJj9gAPAA=`
	b := []byte(msgBody)

	log.Printf("String: %s", msgBody)
	log.Printf("Bytes : %v", b)

	buf := bytes.NewReader(b)
	dec := gob.NewDecoder(buf)

	sd := &dto.SensorMessage{}
	dec.Decode(sd)

	log.Printf("\n\nDecoded String: %v", sd)

}
