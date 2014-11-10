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

func producer(ch chan string, name string, max int) {
	for i := 0; i < max; i++ {
		ch <- fmt.Sprintf("%s:%02d", name, i)
		time.Sleep(5 * time.Second)
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

	err = ch.ExchangeDeclare("chash", "x-consistent-hash", false, false, false, false, nil)
	handleError(err, "Unable to declare exchange")

	msgCh := make(chan string)

	categories := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m"}
	for _, cat := range categories {
		go producer(msgCh, cat, 20)
	}

	counter := len(categories)

	for msg := <-msgCh; ; msg = <-msgCh {
		if "done" == msg {
			counter = counter - 1
		} else {
			err = ch.Publish("chash", msg[0:1], false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg)})
			handleError(err, "unable to send message")

		}
		if counter == 0 {
			break
		}
	}
}
