package amqp

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type (
	Messenger struct {
		channel *amqp.Channel
		queue   *amqp.Queue
	}
	Message struct {
		Username string `json:"username"`
		Message  string `json:"message"`
	}
)

func NewMessenger(dialInfo string) *Messenger {
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
	}
}
func (m *Messenger) SubmitToMessageBroker(message *Message) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	if err := m.channel.Publish(
		"",
		m.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
