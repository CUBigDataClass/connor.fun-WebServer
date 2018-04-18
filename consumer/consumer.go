package consumer

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/tidwall/gjson"
)

type Consumer struct {
	Data       map[string]string
	stream     chan string
	brokerIP   string
	brokerPort string
}

func NewConsumer(stream chan string, brokerIP string, brokerPort string) *Consumer {
	return &Consumer{
		Data:       make(map[string]string),
		stream:     stream,
		brokerIP:   brokerIP,
		brokerPort: brokerPort,
	}
}

func (cons *Consumer) StartConsumer() {
	var mutex sync.Mutex

	topics := []string{"test"}
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":    cons.brokerIP + ":" + cons.brokerPort,
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
				go cons.handleData(string(e.Value), &mutex)
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

func (cons *Consumer) handleData(finalData string, mutex *sync.Mutex) {
	name := gjson.Get(finalData, "name").Str

	mutex.Lock()
	cons.Data[name] = finalData
	mutex.Unlock()

	cons.stream <- finalData
}
