package amqp

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
	"github.com/yhung-mea7/PRO290-twitter-clone/blob/main/notifications-service/data"
)

type Messenger struct {
	channel *amqp.Channel
	queue   *amqp.Queue
	repo    *data.NotificationRepo
}

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func NewMessager(dialInfo string, repo *data.NotificationRepo) *Messenger {
	conn, err := amqp.Dial(dialInfo)
	if err != nil {
		panic(err)
	}
	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	q, err := channel.QueueDeclare(
		"notifications-queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	return &Messenger{
		channel: channel,
		queue:   &q,
		repo:    repo,
	}
}

func (r *Messenger) Consume() {
	fmt.Println("CONSUME")
	msgs, err := r.channel.Consume(
		r.queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	go func() {
		for d := range msgs {
			messageInfo := Message{}
			if err := json.Unmarshal(d.Body, &messageInfo); err != nil {
				fmt.Println(err)
			}
			nc, err := r.repo.RetrieveNotifications(messageInfo.Username)
			if err != nil {
				fmt.Println(err)
			}
			nc.Notification = append(nc.Notification, &messageInfo.Message)
			err = r.repo.SaveNotifications(nc)
			if err != nil {
				fmt.Println(err)
			}
			d.Ack(false)
		}
	}()
}
