package goamqptasks

import (
	"errors"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) error {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		return errors.New(fmt.Sprintf("%s: %s", msg, err))
	}
	return nil
}

type Task struct {
	conn *amqp.Connection
}

func NewTask(uri string) (*Task, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s", uri))
	failOnError(err, "Failed to connect to RabbitMQ")
	return &Task{
		conn,
	}, err
}

func (t *Task) Publish(taskKey, body string) error {
	ch, err := t.conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		fmt.Sprintf("task-%s", taskKey), // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return failOnError(err, "Failed to declare a queue")
	}
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	if err != nil {
		return failOnError(err, "Failed to publish a message")
	}
	return nil
}

func (t *Task) On(taskKey string, handleFunc func(action string, deliveries <-chan amqp.Delivery)) {
	uri := fmt.Sprintf("task-%s", taskKey)
	ch, err := t.conn.Channel()
	failOnError(err, "Failed to open a channel")
	q, err := ch.QueueDeclare(
		uri,   // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	log.Printf("task listen: channel [action='%s'])", uri)
	handleFunc(uri, msgs)
}
