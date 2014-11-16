package main

import (
	"flag"
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

func getQueueName() string {
	if flag.NArg() == 0 {
		panic("Unable to find queue name")
	}
	return flag.Arg(0)
}

func main() {
	log.Printf("Stopping queue on")
	url, err := common.GetRabbitMQ()
	handleError(err, "Unable to find RabbitMQ")

	conn, err := amqp.Dial(url)
	handleError(err, "Unable to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	handleError(err, "Unable to open a channel")
	defer ch.Close()

	const exchangeName = "chash"
	err = ch.ExchangeDeclare(exchangeName, "x-consistent-hash", false, false, false, false, nil)
	handleError(err, "Unable to declare exchange")

	err = ch.QueueUnbind(getQueueName(), "100", exchangeName, nil)
	handleError(err, "Failed to unbind consumer")
}
