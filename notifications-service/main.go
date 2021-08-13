package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/yhung-mea7/PRO290-twitter-clone/blob/main/notifications-service/amqp"
	"github.com/yhung-mea7/PRO290-twitter-clone/blob/main/notifications-service/data"
	"github.com/yhung-mea7/PRO290-twitter-clone/blob/main/notifications-service/handlers"
	"github.com/yhung-mea7/PRO290-twitter-clone/blob/main/notifications-service/register"
)

type (
	userInformation struct {
		Username string `json:"username"`
		UserType int    `json:"user_type"`
	}
)

func main() {
	// sm := mux.NewRouter()
	logger := log.New(os.Stdout, "notification-service", log.LstdFlags)

	consulClient := register.NewConsulClient("notifications-service")
	consulClient.RegisterService()
	defer consulClient.DeregisterService()
	redisCli := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: "",
		DB:       0,
	})
	repo := data.NewNotificationRepo(redisCli)
	nh := handlers.NewNotificationHandler(repo, logger, consulClient, amqp.NewMessager(os.Getenv("RABBIT_CONN"), repo))
	defer redisCli.Close()

	// server := socketio.NewServer(nil)

	// server.OnConnect("/", func(s socketio.Conn) error {
	// 	s.SetContext("")
	// 	log.Println("connected:", s.ID())
	// 	s.Emit("authorization")
	// 	return nil
	// })

	// server.OnEvent("/", "request", func(s socketio.Conn, msg string) {
	// 	resp, err := nh.SendNewRequest("users-service", "GET", "", map[string]string{"Authorization": msg})
	// 	if err != nil {
	// 		logger.Println(err)
	// 	}
	// 	userInfo := userInformation{}
	// 	if err := data.FromJSON(&userInfo, resp.Body); err != nil {
	// 		logger.Println(err)
	// 	}
	// 	resp.Body.Close()
	// 	nc, err := repo.RetrieveNotifications(userInfo.Username)
	// 	if err != nil {
	// 		logger.Println(err)
	// 	}
	// 	notes := nc.Notification
	// 	for _, v := range notes {

	// 		s.Emit("notification", v)
	// 		// notes = notes[1:]
	// 	}
	// 	nc.Notification = notes
	// 	err = repo.SaveNotifications(nc)
	// 	if err != nil {
	// 		logger.Println(err)
	// 	}

	// })

	// server.OnError("/", func(s socketio.Conn, e error) {
	// 	log.Println("meet error:", e)
	// })

	// server.OnDisconnect("/", func(s socketio.Conn, reason string) {
	// 	log.Println("closed", reason)
	// })
	// go func() {
	// 	if err := server.Serve(); err != nil {
	// 		log.Fatalf("socketio listen error: %s\n", err)
	// 	}
	// }()
	// defer server.Close()

	// http.Handle("/", server)
	http.Handle("/healthcheck", nh.HealthCheck())
	logger.Fatal(http.ListenAndServe(os.Getenv("PORT"), nil))
}
