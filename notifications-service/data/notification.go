package data

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

type (
	NotificationCenter struct {
		Username     string    `json:"username"`
		Notification []*string `json:"notifications"`
	}

	NotificationRepo struct {
		db *redis.Client
	}
)

func NewNotificationRepo(cache *redis.Client) *NotificationRepo {
	return &NotificationRepo{
		db: cache,
	}
}

func (nr *NotificationRepo) RetrieveNotifications(username string) (*NotificationCenter, error) {
	val, err := nr.db.Get(username).Result()
	if err == redis.Nil || val == "" {
		fmt.Println("checkpoint")
		newKey := &NotificationCenter{Username: username, Notification: []*string{}}
		return newKey, nr.SaveNotifications(newKey)
	}
	// fmt.Println("checkpoint2")
	outNote := []*string{}
	err = json.Unmarshal([]byte(val), &outNote)
	return &NotificationCenter{username, outNote}, err
}

func (nr *NotificationRepo) SaveNotifications(noteCenter *NotificationCenter) error {
	json, err := json.Marshal(noteCenter.Notification)
	if err != nil {
		return err
	}
	return nr.db.Set(noteCenter.Username, json, 0).Err()
}
