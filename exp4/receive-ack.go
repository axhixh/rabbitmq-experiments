package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
    "time"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	fmt.Println("Receiving message")
	conn, err := amqp.Dial("amqp://admin:VgkDl7PM5DzY@172.17.0.3:5672/")
	handleError(err, "Unable to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	handleError(err, "Unable to open a channel")
	defer ch.Close()

    err = ch.ExchangeDeclare("chash-one", "x-consistent-hash", true, false, false, false, nil)
    handleError(err, "Unable to declare exchange")

    err = ch.Qos(1,0, true)
    handleError(err, "Unable to set prefetch count")

	q, err := ch.QueueDeclare("", true, false, true, false, nil)
	handleError(err, "Unable to create queue")

    err = ch.QueueBind(q.Name, "100", "chash-one", false, nil)
	msg, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	handleError(err, "Failed to register consumer")

	forever := make(chan bool)

	go func() {
		for d := range msg {
			log.Printf("%s: %s", q.Name, d.Body)
            time.Sleep(2 * time.Second)
            d.Ack(false)
		}
	}()

	log.Printf("[*] Waiting. Press CTRL+C to exit")
	<-forever
}
