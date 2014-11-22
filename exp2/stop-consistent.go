package main

import (
	"flag"
	"github.com/axhixh/rabbitmq-experiments/common"
	"github.com/streadway/amqp"
	"log"
)

func getQueueName() string {
	if flag.NArg() == 0 {
		panic("Unable to find queue name")
	}
	return flag.Arg(0)
}

func main() {
	log.Printf("Stopping queue on")
	url, err := common.GetRabbitMQ()
	common.HandleError(err, "Unable to find RabbitMQ")

	conn, err := amqp.Dial(url)
	common.HandleError(err, "Unable to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	common.HandleError(err, "Unable to open a channel")
	defer ch.Close()

	const exchangeName = "chash"
	err = ch.ExchangeDeclare(exchangeName, "x-consistent-hash", false, false, false, false, nil)
	common.HandleError(err, "Unable to declare exchange")

	err = ch.QueueUnbind(getQueueName(), "100", exchangeName, nil)
	common.HandleError(err, "Failed to unbind consumer")
}
