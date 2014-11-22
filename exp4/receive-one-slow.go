package main

import (
	"log"
	"time"

	"github.com/axhixh/rabbitmq-experiments/common"
	"github.com/streadway/amqp"
)

func main() {
	log.Println("Receiving message")
	conn, err := amqp.Dial("amqp://admin:VgkDl7PM5DzY@172.17.0.3:5672/")
	common.HandleError(err, "Unable to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	common.HandleError(err, "Unable to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare("chash-one", "x-consistent-hash", true, false, false, false, nil)
	common.HandleError(err, "Unable to declare exchange")

	err = ch.Qos(1, 0, true)
	common.HandleError(err, "Unable to set prefetch")

	prop := amqp.Table{"x-max-length": int64(1)}
	q, err := ch.QueueDeclare("", true, false, true, false, prop)
	common.HandleError(err, "Unable to create queue")

	err = ch.QueueBind(q.Name, "100", "chash-one", false, nil)
	msg, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	common.HandleError(err, "Failed to register consumer")

	forever := make(chan bool)

	go func() {
		for d := range msg {
			log.Printf("%s: %s", q.Name, d.Body)
			time.Sleep(3 * time.Second)
		}
	}()

	log.Printf("[*] Waiting. Press CTRL+C to exit")
	<-forever
}
