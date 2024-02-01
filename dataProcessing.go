package main

import (
	"github.com/gorilla/websocket"
	"time"
)

// Message represents the overall structure of an incoming websocket message.
type Message struct {
	Data []PriceData `json:"data"`
	Type string      `json:"type"`
}

// PriceData represents the necessary fields from the websocket message.
type PriceData struct {
	P float64 `json:"p"`
	T float64 `json:"t"`
}

type AverageData struct {
	P float64
	T float64
}

// DataGenerator gets data from webscoket and passes to the read channel priceDataChan
func DataGenerator(w *websocket.Conn, priceDataChan chan<- PriceData) {
	for {
		var msg Message
		err := w.ReadJSON(&msg)
		if err != nil {
			panic(err)
		}

		if msg.Type == "trade" {
			for _, data := range msg.Data {
				priceDataChan <- data
			}
		}
	}
}

func SimpleMovingAverage(priceDataChan <-chan PriceData, averageDataChan chan<- AverageData, windowSize int) {
	// variables needed for calculation of SMA
	count := 0
	mean := 0.0
	differential := 0.0
	data := make([]float64, 0, windowSize)

	for elem := range priceDataChan {
		now := time.Now()
		data = append(data, elem.P)
		if len(data) < windowSize {
			count++
			differential = (elem.P - mean) / float64(count)
			mean += differential
			// store the average and the current unix time
			averageDataChan <- AverageData{
				P: mean,
				T: float64(now.UnixNano() / int64(time.Millisecond)),
			}
		} else {
			removedPrice := data[0]
			data = append(data[1:], elem.P)
			differential = (elem.P - removedPrice) / float64(count)
			mean += differential
			averageDataChan <- AverageData{
				P: mean,
				T: float64(now.UnixNano() / int64(time.Millisecond)),
			}
		}

		//fmt.Printf("Average after processing %v: %v\n", elem.P, mean)

	}
}
