# connor.fun-WebServer

This server is how users get the most recent region data when they load the website.  It uses server-side events to send this information.  Every time there is an update to any region, that new data is sent to all of the clients.

## Deployment

### Requirements

This server requires `librdkafka` to run.  This can be installed by following the instructions on the [official GitHub repo](https://github.com/edenhill/librdkafka).

#### Note for AWS

If you are running the standard Amazon Linux distribution on your EC2 instance, or if you just encounter this error somewhere else, there are two things you must do in order to use this library.  First, ensure that you have `gcc-g++` installed prior to beginning installation.  After installation, you will need to export an environment variable.  Simply add `export LD_LIBRARY_PATH=/usr/local/lib` to your `.bashrc` and you're good to go.

### Go Dependencies

If you have `dep` installed, just navigate to the directory and run `dep ensure`.  Otherwise, just run it and see what breaks.  It'll tell you you're missing packages in case you need them.

### Building/Running

The file you want to `go run` and `go build` is called `main.go`.  It takes parameters as command line arguments.  If you do not include these, you will not get an error so always check that the parameters are correct.

To run, simply do `go run main.go <desired-server-host-port> <kafka-broker-ip-address> <kafka-broker-port>`.  You should then be able to see the output from the Kafka consumer.

### Testing

Unit tests?  Who needs them?  The quick and dirty way to test is to open the console in your favorite non-Microsoft browser and run `var evtSource = new EventSource("http://<server-host-ip-address>:<server-host-port>", {} );`.  Then run `evtSource.addEventListener("message", (e) => console.log(e.data))`.  If you don't see any errors, and you eventually see data, then it's working.