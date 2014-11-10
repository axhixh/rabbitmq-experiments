package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func producer(ch chan string, name string, max int) {
	for i := 0; i < max; i++ {
		ch <- fmt.Sprintf("%s:%d", name, i)
	}
	ch <- "done"
}

func main() {
	fmt.Println("Sending message")
	conn, err := amqp.Dial("amqp://admin:O5Wbth9r3F8R@172.17.0.8:5672/")
	handleError(err, "Unable to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	handleError(err, "Unable to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare("q1", false, false, false, false, nil)
	handleError(err, "Unable to create queue")

	msgCh := make(chan string)
	go producer(msgCh, "a", 10)
	go producer(msgCh, "b", 15)
	go producer(msgCh, "c", 7)

	counter := 3

	for msg := <-msgCh; ; msg = <-msgCh {
		if "done" == msg {
			counter = counter - 1
		} else {
			err = ch.Publish("", q.Name, false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg)})
			handleError(err, "unable to send message")

		}
		if counter == 0 {
			break
		}

	}
}
