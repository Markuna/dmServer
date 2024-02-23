package service

import (
	"context"
	"douyinApi/config"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var Ch *amqp.Channel

func init() {
	conn, err := amqp.Dial(config.Get().Rabbitmq.Addr)
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	err = ch.ExchangeDeclare(
		"callback_msg_exchange", // name
		"direct",                // 类型
		false,                   // 持久化
		false,                   // auto-deleted
		false,                   // internal
		false,                   // no-wait
		nil,                     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	Ch = ch
}

func Send(queueName, bodyStr string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := Ch.PublishWithContext(ctx,
		"callback_msg_exchange", // exchange
		queueName,               // routing key
		false,                   // mandatory
		false,                   // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(bodyStr),
		})
	failOnError(err, "Failed to publish a message")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
