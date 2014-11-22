package main

import (
	"github.com/axhixh/rabbitmq-experiments/common"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	log.Println("Receiving message")
	url, err := common.GetRabbitMQ()
	common.HandleError(err, "Unable to get address of RabbitMQ")
	log.Printf(" with RabbitMQ at %s\n", url)
	conn, err := amqp.Dial(url)
	common.HandleError(err, "Unable to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	common.HandleError(err, "Unable to open a channel")
	defer ch.Close()

	const exchangeName = "chash-x"
	err = ch.ExchangeDeclare(exchangeName, "x-consistent-hash", true, false, false, false, nil)
	common.HandleError(err, "Unable to declare exchange")

	q, err := ch.QueueDeclare("", true, false, true, false, nil)
	common.HandleError(err, "Unable to create queue")

	err = ch.QueueBind(q.Name, "100", exchangeName, false, nil)
	msg, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	common.HandleError(err, "Failed to register consumer")

	forever := make(chan bool)

	go func() {
		for d := range msg {
			log.Printf("%s: %s", q.Name, d.Body)
		}
	}()

	log.Printf("[*] Waiting. Press CTRL+C to exit")
	<-forever
}
