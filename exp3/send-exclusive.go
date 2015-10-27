package main

import (
	"github.com/axhixh/rabbitmq-experiments/common"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"time"
)

func main() {
	log.Printf("Sending message")
	url, err := common.GetRabbitMQ()
	common.HandleError(err, "Unable to get address of RabbitMQ")
	log.Printf(" using RabbitMQ %s \n", url)

	conn, err := amqp.Dial(url)
	common.HandleError(err, "Unable to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	common.HandleError(err, "Unable to open a channel")
	defer ch.Close()

	const exchangeName = "chash-x"
	err = ch.ExchangeDeclare(exchangeName, "x-consistent-hash", true, false, false, false, nil)
	common.HandleError(err, "Unable to declare exchange")

	msgCh := make(chan common.Message)

	r := rand.New(rand.NewSource(time.Now().Unix()))
	categories := []string{"AA", "BB", "CC", "DD", "EE", "FF", "GG", "HH"}
	for i, cat := range categories {
		gen := common.Generator{Key: cat, Color: 41 + i}
		log.Printf("Starting %s\n", gen.Key)
		go gen.Generate(msgCh, 20)
		time.Sleep(time.Duration(r.Intn(200)+r.Intn(200)) * time.Millisecond)
	}

	counter := len(categories)

	for msg := <-msgCh; ; msg = <-msgCh {
		if "done" == msg.Body {
			counter = counter - 1
		} else {
			err = ch.Publish(exchangeName, msg.Key, false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg.Body)})
			common.HandleError(err, "unable to send message")

		}
		if counter == 0 {
			break
		}
	}
}
