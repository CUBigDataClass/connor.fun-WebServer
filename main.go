package main

import (
	"os"

	"github.com/CUBigDataClass/connor.fun-Kafka/server"
)

func main() {
	serv := server.NewServer()
	serv.StartServer(os.Args[1], os.Args[2], os.Args[3])
}
