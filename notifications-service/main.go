package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis"
	"github.com/yhung-mea7/PRO290-twitter-clone/blob/main/notifications-service/amqp"
	"github.com/yhung-mea7/PRO290-twitter-clone/blob/main/notifications-service/data"
	"github.com/yhung-mea7/PRO290-twitter-clone/blob/main/notifications-service/handlers"
	"github.com/yhung-mea7/PRO290-twitter-clone/blob/main/notifications-service/register"
)

func main() {
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

	server := nh.NotificationConnection()
	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer server.Close()

	http.Handle("/socket.io/", server)

	// routes.SetUpRoutes(sm, nh)
	http.Handle("/healthcheck", nh.HealthCheck())
	// http.Handle("/", nh.NotificationConnection())

	serve := http.Server{
		Addr:     os.Getenv("PORT"),
		Handler:  http.DefaultServeMux,
		ErrorLog: logger,
		// ReadTimeout:  100 * time.Second,
		// WriteTimeout: 100 * time.Second,
		// IdleTimeout:  120 * time.Second,
	}

	go func() {
		logger.Printf("Starting server on port: %v \n", serve.Addr)
		err := serve.ListenAndServe()
		if err != nil {
			logger.Printf("Error starting server: %v \n", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	sig := <-c
	logger.Println("Got Signal:", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	serve.Shutdown(ctx)
}
