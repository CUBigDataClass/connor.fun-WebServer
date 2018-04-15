## connor.fun-Kafka

This is the current place development location for kafka consumption/data serving.  Please ignore the commits to master, I forgot to create a new branch and it was too late.

Currently, there is a small self-contained build that you can run.  This can be found in the LocalTest directory. 

### Testing Locally

First, you need to run Kafka locally.  Follow this quickstart guide to see how it is done: https://kafka.apache.org/quickstart

In order to run the local consumer/producer, you must first install `rdkafka`.  Clone this repo and follow the instructions in the README: https://github.com/edenhill/librdkafka

The local consumer/producers require the local broker ID to equal 0, for the broker server/port to remain at the default settings, and for the topic to be named "test".

Once Kafka is running with those parameters, `./producer` will randomly select a sample JSON object formatted how the final data will (hopefully) be formatted.  Running `./hackyServer` will start a server at `localhost:8000`.  This will accept websocket connections.  Once a new connection is established, all data that has been collected so far will be sent to the new client.  

**NOTE**: The server does not get any data that is already logged when it starts.  This means that `./producer` needs to be run before making a new websocket connection.  In order to get initial data.