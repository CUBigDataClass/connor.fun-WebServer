package main

import (
	"os"

	"github.com/CUBigDataClass/connor.fun-Kafka/server"
)

func main() {
	serv := server.NewServer()
	serverPort := string(os.Args[1])
	brokerIP := string(os.Args[2])
	brokerPort := string(os.Args[3])
	serv.StartServer(serverPort, brokerIP, brokerPort)
}
