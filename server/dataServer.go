package server

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/CUBigDataClass/connor.fun-Kafka/consumer"
	"github.com/aaronaaeng/event-source-add"
)

type DataServer struct {
	serverPort string
	dataStream chan string
	cons       *consumer.Consumer
	messageID  int
	es         eventsource.EventSource
	addAlert   chan bool
}

func NewServer(serverPort string, brokerIP string, brokerPort string) {
	stream := make(chan string)
	addAlert := make(chan bool)

	cons := consumer.NewConsumer(stream, brokerIP, brokerPort)
	serv := DataServer{
		serverPort: serverPort,
		dataStream: stream,
		cons:       cons,
		messageID:  0,
		es: eventsource.New(
			addAlert,
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

	log.Fatal(http.ListenAndServe(":"+serv.serverPort, nil))

	go func() {
		for {
			select {
			case data := <-serv.dataStream:
				serv.es.SendEventMessage(data, "message", strconv.Itoa(serv.messageID))
				serv.messageID++
			case addedClient := <-serv.addAlert:
				if addedClient {
					for _, data := range serv.cons.Data {
						serv.es.SendEventMessage(data, "message", strconv.Itoa(serv.messageID))
						serv.messageID++
					}
				}
			}
		}
	}()
}
