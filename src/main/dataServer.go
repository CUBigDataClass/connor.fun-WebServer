package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []byte)
var upgrader = websocket.Upgrader{}

// type dataServer struct {
// 	cons *consumer
// }

// func newServer() *dataServer {
// 	return &dataServer{
// 		cons: consumer.newConsumer(),
// 	}
// }

func main() {
	// serv := newServer()

	// go serv.cons.startConsumer()

	// fs := http.FileServer(http.Dir("/home/odin/Go/websocket-chat/public"))
	// http.Handle("/", fs)

	http.HandleFunc("/", handleConnections)

	http.HandleFunc("/ws", handleConnections)

	go handleUpdates()

	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	clients[ws] = true
	// serv.sendCached(ws)
}

// func (serv *dataServer) sendCached(client *websocket.Conn) {
// 	for data := range serv.cons.data {
// 		err := client.WriteJSON(data)
// 		if err != nil {
// 			log.Printf("error: %v", err)
// 			client.Close()
// 			delete(clients, client)
// 		}
// 	}
// }

func handleUpdates() {
	for {
		data := <-broadcast

		for client := range clients {
			err := client.WriteJSON(data)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
