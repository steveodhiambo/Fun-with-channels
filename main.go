package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"os"
)

func main() {
	apiKey := os.Getenv("FINNHUB_TOKEN")

	w, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://ws.finnhub.io?token=%s", apiKey), nil)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	symbols := []string{"BINANCE:BTCUSDT"} // "BINANCE:ETHUSDT", "BINANCE:ADAUSDT"}
	for _, s := range symbols {
		msg, _ := json.Marshal(map[string]interface{}{"type": "subscribe", "symbol": s})

		err = w.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			panic(err)
		}
	}

	// Channel for passing price data
	priceDataChan := make(chan PriceData)
	defer close(priceDataChan)

	// channel for storing average data
	averageDataChan := make(chan AverageData)
	defer close(averageDataChan)

	// DataGenerator
	go DataGenerator(w, priceDataChan)

	// Calculate Simple Moving Average
	go SimpleMovingAverage(priceDataChan, averageDataChan, 60)

	// Create file to store the average data
	averageDataFile, err := os.OpenFile("data/average.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer averageDataFile.Close()

	// Create a writer that writes to file and also prints to stdout
	wrt := io.MultiWriter(os.Stdout, averageDataFile)
	log.SetOutput(wrt)

	for elem := range averageDataChan {
		log.Println("-------------------------")
		log.Printf("Average:  %+v\n", elem.P)
		log.Printf("Time:  %+v\n", elem.T)
		log.Println("-------------------------")
	}

}
