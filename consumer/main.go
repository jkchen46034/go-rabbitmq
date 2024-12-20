package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://localhost:5672")
	Panic(err)
	defer conn.Close()

	channel, err := conn.Channel()
	Panic(err)
	defer channel.Close()

	queueName := "temperature"
	messages, err := channel.Consume(queueName, "", true, false, false, false, nil)
	Panic(err)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case message := <-messages:
			log.Printf("Message: %s\n", message.Body)
		case <-sigchan:
			log.Println("interrupt detected")
			os.Exit(0)
		}
	}
}
