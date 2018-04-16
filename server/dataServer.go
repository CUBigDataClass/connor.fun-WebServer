package server

import (
	"fmt"
	"net/http"

	"github.com/CUBigDataClass/connor.fun-Kafka/consumer"
	"github.com/gorilla/websocket"
)

type DataServer struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	consumer  *consumer.Consumer
	upgrader  websocket.Upgrader
}

func NewServer() *DataServer {
	broadcastStream := make(chan []byte)
	return &DataServer{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: broadcastStream,
		consumer:  consumer.NewConsumer(broadcastStream),
		upgrader:  websocket.Upgrader{},
	}
}

func (serv *DataServer) StartServer(serverPort string, brokerIP string, brokerPort string) {
	go serv.consumer.StartConsumer(brokerIP, brokerPort)

	http.HandleFunc("/", serv.handleConnections)

	http.HandleFunc("/ws", serv.handleConnections)

	go serv.handleUpdates()

	fmt.Println("http server started on :%s", serverPort)
	fmt.Println(":" + serverPort)
	err := http.ListenAndServe(":"+serverPort, nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}

func (serv *DataServer) handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := serv.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}

	serv.clients[ws] = true
	serv.sendCached(ws)
}

func (serv *DataServer) sendCached(client *websocket.Conn) {
	fmt.Println("In sendCached")
	for name := range serv.consumer.Data {
		err := client.WriteJSON(string(serv.consumer.Data[name]))
		if err != nil {
			fmt.Printf("error: %v", err)
			client.Close()
			delete(serv.clients, client)
		}
	}
}

func (serv *DataServer) handleUpdates() {
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
