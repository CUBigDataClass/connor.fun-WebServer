package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/CUBigDataClass/connor.fun-Kafka/consumer"
	"github.com/antage/eventsource"
)

func main() {
	serverPort := string(os.Args[1])
	brokerIP := string(os.Args[2])
	brokerPort := string(os.Args[3])
	dataStream := make(chan string)

	cons := consumer.NewConsumer(dataStream, brokerIP, brokerPort)
	go cons.StartConsumer()

	es := eventsource.New(
		&eventsource.Settings{
			Timeout:        5 * time.Second,
			CloseOnTimeout: false,
			IdleTimeout:    30 * time.Minute,
		},
		func(req *http.Request) [][]byte {
			return [][]byte{
				[]byte("X-Accel-Buffering: no"),
				[]byte("Access-Control-Allow-Origin: *"),
				[]byte("Access-Control-Expose-Headers: *"),
				[]byte("Access-Control-Allow-Credentials: true"),
			}
		})

	defer es.Close()

	http.Handle("/", es)

	go func() {
		messageID := 0
		for {
			select {
			case data := <-dataStream:
				es.SendEventMessage(data, "message", strconv.Itoa(messageID))
				messageID++
			default:
			}
		}
	}()
	log.Fatal(http.ListenAndServe(":"+serverPort, nil))
}
