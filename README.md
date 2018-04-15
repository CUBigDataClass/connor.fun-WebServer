## connor.fun-Kafka

This is the current place development location for kafka consumption/data serving.  Please ignore the commits to master, I forgot to create a new branch and it was too late.

Currently, there is a small self-contained build that you can run.  This can be found in the LocalTest directory. 

### Manual Local Test

First, you need to run Kafka locally.  Follow this quickstart guide to see how it is done: https://kafka.apache.org/quickstart

In order to run the local consumer/producer, you must first install `rdkafka`.  Clone this repo and follow the instructions in the README: https://github.com/edenhill/librdkafka

The local consumer/producers require the local broker ID to equal 0, for the broker server/port to remain at the default settings, and for the topic to be named "test".

Once Kafka is running with those parameters, `./producer` will randomly select a sample JSON object formatted how the final data will be formatted.  Every 10 seconds, it will randomly select one of nine sample data-points and put it in the queue.  Running `./hackyServer` will start a server at `localhost:8000`.  This will accept websocket connections.  Once a new connection is established, all data that has been collected so far will be sent to the new client.  

**NOTE**: The server does not get any data that is already logged when it starts.  This means that `./producer` needs to be run before making a new websocket connection. 

#### Automated Local Test

I have created a simple script called `auto-start.sh`.  This can be found in the `LocalTest` directory.  This directory even contains the Kafka source because that's what I decided to do for you.  `./auto-run.sh` will start the zookeeper and a broker.  `./produer` will start the looping producer.  From here, run `./hackyServer` to start a server.  This shell script is the sketchiest thing I've ever written.  You will likely see error messages saying that Kafka encountered a fatal error and is shutting down.  However, it never does shut down.  There is also a chance that the looping producer never exists.  If that risk is too much for you, you can instead do `./produceOnce`.  Don't ask me why it doesn't just take command-line arguments.
