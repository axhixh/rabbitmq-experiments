package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
    "math/rand"
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
		time.Sleep(2 * time.Second)
	}
	ch <- "done"
}

func main() {
	fmt.Println("Sending message")
	conn, err := amqp.Dial("amqp://admin:VgkDl7PM5DzY@172.17.0.3:5672/")
	handleError(err, "Unable to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	handleError(err, "Unable to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare("chash-x", "x-consistent-hash", true, false, false, false, nil)
	handleError(err, "Unable to declare exchange")

	msgCh := make(chan string)

    r := rand.New(rand.NewSource(time.Now().Unix()))
	categories := []string{"AA", "BB", "CC", "DD", "EE", "FF", "GG", "HH"}
	for i, cat := range categories {
        m := fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", 40 + i, cat)
		go producer(msgCh, m, 50)
		time.Sleep(time.Duration(r.Intn(200) + r.Intn(200)) * time.Millisecond)
	}

	counter := len(categories)

	for msg := <-msgCh; ; msg = <-msgCh {
		if "done" == msg {
			counter = counter - 1
		} else {
			err = ch.Publish("chash-x", msg[7:9], false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg)})
			handleError(err, "unable to send message")

		}
		if counter == 0 {
			break
		}
	}
}
