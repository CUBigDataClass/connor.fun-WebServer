package server

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/CUBigDataClass/connor.fun-Kafka/consumer"
	"github.com/aaronaaeng/eventsource"
)

type DataServer struct {
	serverPort string
	dataStream chan string
	cons       *consumer.Consumer
	messageID  int
	es         eventsource.EventSource
}

func NewServer(serverPort string, brokerIP string, brokerPort string) {
	var stream chan string

	cons := consumer.NewConsumer(stream, brokerIP, brokerPort)
	serv := DataServer{
		serverPort: serverPort,
		dataStream: stream,
		cons:       cons,
		messageID:  0,
		es: eventsource.New(
			cons,
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
			}),
	}
	go serv.cons.StartConsumer()

	serv.startServer()
}

func (serv *DataServer) startServer() {
	defer serv.es.Close()
	http.Handle("/", serv.es)
	go serv.sendUpdates()

	log.Fatal(http.ListenAndServe(":"+serv.serverPort, nil))
}

func (serv *DataServer) sendUpdates() {
	for {
		select {
		case data := <-serv.dataStream:
			serv.es.SendEventMessage(data, "message", strconv.Itoa(serv.messageID))
			serv.messageID++
		}
	}
}
