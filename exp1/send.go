package main

import (
	"fmt"
	"github.com/axhixh/rabbitmq-experiments/common"
	"github.com/streadway/amqp"
	"log"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	log.Printf("Sending message")
	url, err := common.GetRabbitMQ()
	handleError(err, "Unable to get RabbitMQ")

	conn, err := amqp.Dial(url)
	handleError(err, "Unable to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	handleError(err, "Unable to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare("q1", false, false, false, false, nil)
	handleError(err, "Unable to create queue")

	generators := []common.Generator{
		common.Generator{Key: "AA", Color: 41},
		common.Generator{Key: "BB", Color: 42},
		common.Generator{Key: "CC", Color: 43}}

	msgCh := make(chan common.Message)
	for i := range generators {
		log.Printf("Starting %s", generators[i].Key)
		go generators[i].Generate(msgCh, 6)
	}

	counter := len(generators)

	for msg := <-msgCh; ; msg = <-msgCh {
		if "done" == msg.Body {
			counter--
			log.Printf("finished %s", msg.Key)
		} else {
			log.Printf("sending %s", msg.Body)
			err = ch.Publish("", q.Name, false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg.Body)})
			handleError(err, "unable to send message")

		}
		if counter == 0 {
			break
		}

	}
}
