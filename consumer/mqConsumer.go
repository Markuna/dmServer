package consumer

import (
	"context"
	"douyinApi/config"
	ct "douyinApi/consumer/convert"
	"douyinApi/consumer/proto"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	LiveComment = "live_comment" // 评论
	LiveGift    = "live_gift"    // 礼物
	LiveLike    = "live_like"    // 点赞
)

var Ch *amqp.Channel

func init() {
	conn, err := amqp.Dial(config.Get().Rabbitmq.Addr)
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	Ch = ch
}

func BuildConsumer(ctx context.Context, roomId string) <-chan *proto.PushMsg {

	out := make(chan *proto.PushMsg)
	msgs1, _ := tryToBindConsumer(LiveComment, roomId)
	msgs2, _ := tryToBindConsumer(LiveGift, roomId)
	msgs3, _ := tryToBindConsumer(LiveLike, roomId)

	go func() {
		for {
			select {
			case d := <-msgs1:
				out <- ct.ConvertPushMsg(LiveComment, d.Body)
			case d := <-msgs2:
				out <- ct.ConvertPushMsg(LiveGift, d.Body)
			case d := <-msgs3:
				out <- ct.ConvertPushMsg(LiveLike, d.Body)
			case <-ctx.Done():
				Ch.Cancel("douyin."+LiveComment+"."+roomId+".consumer", false)
				Ch.Cancel("douyin."+LiveGift+"."+roomId+".consumer", false)
				Ch.Cancel("douyin."+LiveLike+"."+roomId+".consumer", false)
				return
			}
		}
	}()

	return out
}
func tryToBindConsumer(msgType, roomId string) (<-chan amqp.Delivery, error) {

	q, err := Ch.QueueDeclare(
		"douyin."+msgType+"."+roomId, // name
		false,                        // durable
		true,                         // delete when unused
		true,                         // exclusive
		true,                         // no-wait
		nil,                          // arguments
	)
	failOnError(err, "Failed to declare a queue")
	err = Ch.QueueBind(
		q.Name,                  // queue name
		q.Name,                  // routing key
		"callback_msg_exchange", // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")
	msgs, err := Ch.Consume(
		q.Name,             // queue
		q.Name+".consumer", // consumer
		true,               // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	failOnError(err, "Failed to consume a queue")
	log.Println("build ", q.Name)
	return msgs, nil
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
