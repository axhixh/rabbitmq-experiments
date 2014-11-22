package main

import (
	"flag"
	"github.com/axhixh/rabbitmq-experiments/common"
	"github.com/streadway/amqp"
	"log"
)

func getQueueName() string {
	if flag.NArg() == 0 {
		panic("Please provide a queue name")
	}
	return flag.Arg(0)
}

func main() {
	log.Printf("Receiving message ")
	url, err := common.GetRabbitMQ()
	common.HandleError(err, "Unable to get RabbitMQ")

	queueName := getQueueName()

	conn, err := amqp.Dial(url)
	common.HandleError(err, "Unable to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	common.HandleError(err, "Unable to open a channel")
	defer ch.Close()

	const exchangeName = "chash"
	err = ch.ExchangeDeclare(exchangeName, "x-consistent-hash", false, false, false, false, nil)
	common.HandleError(err, "Unable to declare exchange")

	q, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
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
