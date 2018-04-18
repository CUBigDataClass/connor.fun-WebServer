package main

import (
	"os"

	"github.com/CUBigDataClass/connor.fun-Kafka/server"
)

func main() {
	serverPort := string(os.Args[1])
	brokerIP := string(os.Args[2])
	brokerPort := string(os.Args[3])
	server.NewServer(serverPort, brokerIP, brokerPort)
}
