#RabbitMQ Experiments

Project to experiment with consistent hash exchange in RabbitMQ

In each experiment the invidiual _go_ source file is an 
independent binary. Please build them separately using the 
command:

    go build <file.go>

This means this project isn't *go-gettable*


#Experiment 1: Basic Queue
Generates messages with different keys and sequence number.
Sends to a generic queue which can have multiple consumers.
The messages are consumed by consumers in round-robin.

#Experiment 2: Consistent Hash
Use consistent hash exchange to distribute messages with
same key to same consumer. Messages are processed in order.

Stop unbinds the queue from the exchange.

#Experiment 3: Exclusive Queue
Consistent hash and exclusive queue.

Shows how when a slow consumer is involved crashing consumer
results in lost message.

#Experiment 4: Limiting Queue Size
