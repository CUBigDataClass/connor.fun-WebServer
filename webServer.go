package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/CUBigDataClass/connor.fun-WebServer/consumer"
	"github.com/antage/eventsource"
	"github.com/tidwall/gjson"
)

type WebServer struct {
	serverPort string
	dataStream chan string
	locations  []byte
	data       map[string]string
	dataMut    *sync.Mutex
	cons       *consumer.Consumer
	es         eventsource.EventSource
}

func main() {
	parameters := "<desired-server-host-port> <kafka-broker-ip-address> <kafka-broker-port>"
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s %s\n",
			os.Args[0], parameters)
		os.Exit(1)
	}

	serverPort := string(os.Args[1])
	brokerIP := string(os.Args[2])
	brokerPort := string(os.Args[3])

	dataStream := make(chan string)
	data := make(map[string]string)
	var mutex sync.Mutex

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
		locations:  loadLocations("./locations.json"),
		data:       data,
		dataMut:    &mutex,
		cons:       cons,
		es:         es,
	}

	defer es.Close()

	http.HandleFunc("/current", serv.getCurrentData)
	http.HandleFunc("/locations", serv.getCurrentLocations)
	http.Handle("/", es)

	go func() {
		messageID := 0
		for {
			select {
			case data := <-dataStream:
				es.SendEventMessage(data, "message", strconv.Itoa(messageID))

				serv.dataMut.Lock()
				key := gjson.Get(data, "ID").Str
				serv.data[key] = data
				serv.dataMut.Unlock()

				messageID++
			}
		}
	}()
	log.Fatal(http.ListenAndServe(":"+serverPort, nil))
}

func loadLocations(file string) []byte {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return raw
}

func (serv *WebServer) getCurrentData(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/current" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		response := "["
		serv.dataMut.Lock()
		for _, data := range serv.data {
			response += data + ","
		}
		response = strings.TrimSuffix(response, ",")
		serv.dataMut.Unlock()
		response += "]"

		sendResponse(w, r, response)

	default:
		fmt.Fprintf(w, "Bad GET request")
	}
}

func (serv *WebServer) getCurrentLocations(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/locations" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		response := string(serv.locations)
		sendResponse(w, r, response)

	default:
		fmt.Fprintf(w, "Bad GET request")
	}
}

func sendResponse(w http.ResponseWriter, r *http.Request, response string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}
