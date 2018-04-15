package main

import (
	"fmt"
	"net/http"

	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/tidwall/gjson"
)

type dataServer struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	consumer  *Consumer
	upgrader  websocket.Upgrader
}

func newServer() *dataServer {
	broadcastStream := make(chan []byte)
	return &dataServer{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: broadcastStream,
		consumer:  NewConsumer(broadcastStream),
		upgrader:  websocket.Upgrader{},
	}
}

func main() {
	serv := newServer()

	go serv.consumer.StartConsumer()

	http.HandleFunc("/", serv.handleConnections)

	http.HandleFunc("/ws", serv.handleConnections)

	go serv.handleUpdates()

	fmt.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}

func (serv *dataServer) handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := serv.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}

	serv.clients[ws] = true
	serv.sendCached(ws)
}

func (serv *dataServer) sendCached(client *websocket.Conn) {
	fmt.Println("In sendCached")
	for data := range serv.consumer.data {
		err := client.WriteJSON(string(data))
		if err != nil {
			fmt.Printf("error: %v", err)
			client.Close()
			delete(serv.clients, client)
		}
	}
}

func (serv *dataServer) handleUpdates() {
	for {
		data := <-serv.broadcast

		for client := range serv.clients {
			err := client.WriteJSON(string(data))
			if err != nil {
				fmt.Printf("error: %v", err)
				client.Close()
				delete(serv.clients, client)
			}
		}
	}
}

type Consumer struct {
	data   map[string][]byte
	stream chan []byte
}

func NewConsumer(stream chan []byte) *Consumer {
	return &Consumer{
		data:   make(map[string][]byte),
		stream: stream,
	}
}

func (cons *Consumer) StartConsumer() {
	topics := []string{"test"}
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":    "localhost:9092",
		"group.id":             "cons",
		"session.timeout.ms":   6000,
		"default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"}})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Consumer %v\n", c)

	err = c.SubscribeTopics(topics, nil)

	run := true

	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			os.Exit(1)
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
				go cons.handleData(e.Value)
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				run = false
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	c.Close()
	fmt.Printf("Closing consumer\n")
}

func (cons *Consumer) handleData(finalData []byte) {
	name := gjson.Get(string(finalData), "name").Str
	cons.data[name] = finalData

	cons.stream <- finalData
}
