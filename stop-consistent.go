package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
    "os"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	fmt.Println("Stopping queue")
	conn, err := amqp.Dial("amqp://admin:O5Wbth9r3F8R@172.17.0.8:5672/")
	handleError(err, "Unable to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	handleError(err, "Unable to open a channel")
	defer ch.Close()

    err = ch.ExchangeDeclare("chash", "x-consistent-hash", false, false, false, false, nil)
    handleError(err, "Unable to declare exchange")

    err = ch.QueueUnbind(os.Args[1], "100", "chash", nil)
	handleError(err, "Failed to unbind consumer")
}
