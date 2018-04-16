// Example function-based Apache Kafka producer
package main

/**
 * Copyright 2016 Confluent Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var jsonHell = [9]string{`{ "name": "New York", "ID": "NY0", "sentiment": -0.0030300168, "tid":"984908849458110464", "centerLat":40.7128, "centerLon":-74.006, "weather": { "main":"Clear", "description":"clear sky", "temp": 84.00 } }`,
	`{ "name": "New York", "ID": "NY0", "sentiment": -0.0030257932, "tid":"984909342137815040", "centerLat":40.7128, "centerLon":-74.006, "weather": { "main":"Clear", "description":"clear sky", "temp": 82.00 } }`,
	`{ "name": "New York", "ID": "NY0", "sentiment": -0.0031428882, "tid":"984909850076499970", "centerLat":40.7128, "centerLon":-74.006, "weather": { "main":"Clear", "description":"clear sky", "temp": 81.75 } }`,
	`{ "name":"Los Angeles", "ID": "LA0", "sentiment":-0.006765286, "tid":"984908841497149440", "centerLat":34.0522, "centerLon":-118.2437, "weather": { "main":"Dust", "description":"dust", "temp": 76.75 } }`,
	`{ "name":"Los Angeles", "ID": "LA0", "sentiment":-0.0068500843, "tid":"984909365500174336", "centerLat":34.0522, "centerLon":-118.2437, "weather": { "main":"Dust", "description":"dust", "temp": 75.75 } }`,
	`{ "name":"Los Angeles", "ID": "LA0", "sentiment":-0.00699275, "tid":"984909860306280449", "centerLat":34.0522, "centerLon":-118.2437, "weather": { "main":"Dust", "description":"dust", "temp": 71.75 } }`,
	`{ "name":"Seattle", "ID": "SEA0", "sentiment":-0.0026830027, "tid":"984908889740034049", "centerLat":47.6062, "centerLon":-122.3321, "weather": { "main":"Rain", "description":"light rain", "temp": 49.60 } }`,
	`{ "name":"Seattle", "ID": "SEA0", "sentiment":-0.0026437445, "tid":"984909868053118976", "centerLat":47.6062, "centerLon":-122.3321, "weather": { "main":"Rain", "description":"light rain", "temp": 48.75 } }`,
	`{ "name":"Seattle", "ID": "SEA0", "sentiment":-0.00699275, "tid":"984909860306280449", "centerLat":47.6062, "centerLon":-122.3321, "weather": { "main":"Drizzle", "description":"light intensity drizzle", "temp": 49.8 } } ]}`}

func main() {
	broker := 0
	topic := "test"
	p, _ := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	sigchan := make(chan os.Signal, 1)

	for {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			os.Exit(1)
		default:
			sendMessage(p, topic)
			time.Sleep(10 * time.Second)
		}
	}

	// Optional delivery channel, if not specified the Producer object's
	// .Events channel is used.
}

func sendMessage(p *kafka.Producer, topic string) {

	deliveryChan := make(chan kafka.Event)

	value := getFakeData()
	err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(value),
		Headers:        []kafka.Header{{"myTestHeader", []byte("header values are binary")}},
	}, deliveryChan)

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}
	close(deliveryChan)
}

func getFakeData() []byte {
	rand.Seed(time.Now().Unix())
	n := rand.Int() % 9
	return []byte(jsonHell[n])
}
