package main

import (
	"fmt"
	"github.com/axhixh/rabbitmq-experiments/common"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	log.Printf("Receiving messages from: ")
	url, err := common.GetRabbitMQ()
	common.HandleError(err, "Unable to get URL for RabbitMQ")

	fmt.Printf("%s\n", url)
	conn, err := amqp.Dial(url)
	common.HandleError(err, "Unable to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	common.HandleError(err, "Unable to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare("q1", false, false, false, false, nil)
	common.HandleError(err, "Unable to create queue")

	msg, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	common.HandleError(err, "Failed to register consumer")

	forever := make(chan bool)

	go func() {
		for d := range msg {
			log.Printf("Msg: %s", d.Body)
		}
	}()

	log.Printf("[*] Waiting. Press CTRL+C to exit")
	<-forever
}
