package common

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)

func Receive(msgCh <-chan amqp.Delivery, queueName string, delay time.Duration) {
	forever := make(chan bool)

	go func() {
		for d := range msgCh {
			log.Printf("%s: %s", queueName, d.Body)
			if delay > 0 {
				time.Sleep(delay * time.Second)
			}
		}
	}()

	<-forever
}
