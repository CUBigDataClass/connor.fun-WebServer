#!/bin/bash
# call the other script using absolute path
./kafka/bin/zookeeper-server-start.sh ./kafka/config/zookeeper.properties & 

./kafka/bin/kafka-server-start.sh ./kafka/config/server.properties &

./kafka/bin/kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic test &

./producer 

