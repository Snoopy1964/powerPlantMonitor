package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/Snoopy1964/powerPlantMonitor/distributed/dto"
)

func main() {
	/*
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

	*/
	fmt.Println("... Started ....")
	buf := new(bytes.Buffer)

	reading := dto.SensorMessage{
		Name:  "hallo",
		Value: 3.1415,
		Tst:   time.Now(),
	}

	buf.Reset()
	enc := gob.NewEncoder(buf)
	enc.Encode(reading)

	fmt.Printf("Bytes:  %v\n", buf.Bytes())
	fmt.Printf("String: %s\n", string(buf.Bytes()))

	// ------------------- decoding
	r := bytes.NewReader(buf.Bytes())
	d := gob.NewDecoder(r)
	sd := new(dto.SensorMessage)
	d.Decode(sd)
	fmt.Printf("Decoded: %v\n", sd)

}
