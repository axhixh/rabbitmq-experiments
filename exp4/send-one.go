package main

import (
	"fmt"
	"github.com/axhixh/rabbitmq-experiments/common"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"time"
)

func producer(ch chan string, name string, max int) {
	for i := 0; i < max; i++ {
		ch <- fmt.Sprintf("%s:%02d", name, i)
		time.Sleep(2 * time.Second)
	}
	ch <- "done"
}

func main() {
	log.Println("Sending message")
	conn, err := amqp.Dial("amqp://admin:VgkDl7PM5DzY@172.17.0.3:5672/")
	common.HandleError(err, "Unable to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	common.HandleError(err, "Unable to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare("chash-one", "x-consistent-hash", true, false, false, false, nil)
	common.HandleError(err, "Unable to declare exchange")

	msgCh := make(chan common.Message)

	r := rand.New(rand.NewSource(time.Now().Unix()))
	categories := []string{"AA", "BB", "CC", "DD", "EE", "FF", "GG", "HH"}
	for i, cat := range categories {
		gen := common.Generator{Key: cat, Color: 41 + i}
		log.Printf("Starting %s", gen.Key)
		go gen.Generate(msgCh, 20)
		time.Sleep(time.Duration(r.Intn(200)+r.Intn(200)) * time.Millisecond)
	}

	counter := len(categories)

	for msg := <-msgCh; ; msg = <-msgCh {
		if "done" == msg.Body {
			counter = counter - 1
		} else {
			err = ch.Publish("chash-one", msg.Key, false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg.Body)})
			common.HandleError(err, "unable to send message")

		}
		if counter == 0 {
			break
		}
	}
}
