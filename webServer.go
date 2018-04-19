package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/CUBigDataClass/connor.fun-Kafka/consumer"
	"github.com/antage/eventsource"
)

type WebServer struct {
	serverPort string
	dataStream chan string
	cons       *consumer.Consumer
	es         eventsource.EventSource
}

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

	serv := &WebServer{
		serverPort: serverPort,
		dataStream: dataStream,
		cons:       cons,
		es:         es,
	}

	defer es.Close()

	http.HandleFunc("/current", serv.getCurrent)
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

func (serv *WebServer) getCurrent(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/current" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		response := "["

		for _, data := range serv.cons.Data {
			response += data + ", "
		}
		response += "]"
		if err := json.NewEncoder(w).Encode(response); err != nil {
			panic(err)
		}
	default:
		fmt.Fprintf(w, "Bad GET request")
	}
}
