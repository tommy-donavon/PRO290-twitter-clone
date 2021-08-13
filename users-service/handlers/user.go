package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/users-service/amqp"
	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/users-service/auth"
	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/users-service/data"
)

type (
	UserHandler struct {
		repo      *data.UserRepo
		log       *log.Logger
		jwt       *auth.JwtWrapper
		messenger *amqp.Messenger
	}
	generalMessage struct {
		Message interface{} `json:"message"`
	}
	userInformation struct {
		Username string `json:"username"`
		UserType int    `json:"user_type"`
	}
	keyValue struct{}
)

func NewUserHandler(repo *data.UserRepo, log *log.Logger, key string, messanger *amqp.Messenger) *UserHandler {
	return &UserHandler{
		repo: repo,
		log:  log,
		jwt: &auth.JwtWrapper{
			SecretKey:       key,
			Issuer:          "users-service",
			ExpirationHours: 24,
		},
		messenger: messanger,
	}
}

func getUserName(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["username"]
}
