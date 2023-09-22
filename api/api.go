package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Todo struct {
	Title       string `json:"title"`
	IsCompleted bool   `json:"isCompleted"`
}

func main() {
	conn, err := amqp.Dial("amqp://test:test@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"todo",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	m := &Todo{Title: "FirstTodo", IsCompleted: false}

	j, err := json.Marshal(m)

	if err != nil {
		return
	}

	err = ch.PublishWithContext(ctx,
		"todo",        // exchange
		"todo.create", // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(j),
		})
	failOnError(err, "Failed to publish a message")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
