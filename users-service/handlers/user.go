package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/users-service/auth"
	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/users-service/data"
)

type (
	UserHandler struct {
		repo *data.UserRepo
		log  *log.Logger
		jwt  *auth.JwtWrapper
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

func NewUserHandler(repo *data.UserRepo, log *log.Logger, key string) *UserHandler {
	return &UserHandler{
		repo: repo,
		log:  log,
		jwt: &auth.JwtWrapper{
			SecretKey:       key,
			Issuer:          "users-service",
			ExpirationHours: 24,
		},
	}
}

func (uh *UserHandler) CreateUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		uh.log.Println("POST CREATE USER")
		user := r.Context().Value(keyValue{}).(data.User)
		if err := uh.repo.CreateUser(&user); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			uh.log.Println(err)
			data.ToJSON(&generalMessage{err.Error()}, rw)
			return
		}
		rw.WriteHeader(http.StatusNoContent)

	}
}

func (uh *UserHandler) DeleteUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userInformation := r.Context().Value(keyValue{}).(userInformation)
		if err := uh.repo.DeleteUser(userInformation.Username); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMessage{err.Error()}, rw)
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}

func (uh *UserHandler) LoginUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		uh.log.Println("POST LOGIN")
		login := r.Context().Value(keyValue{}).(data.Login)
		user, err := uh.repo.GetUser(login.Username)
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			data.ToJSON(&generalMessage{"Invalid Login information"}, rw)
			return
		}
		if err := user.CheckPassword(login.Password); err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			data.ToJSON(&generalMessage{"Invalid Login information"}, rw)
			return
		}
		token, err := uh.jwt.CreateJwToken(user.Username, user.UserType)
		if err != nil {
			uh.log.Println(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMessage{"Failed to sign token"}, rw)
			return
		}
		data.ToJSON(
			struct {
				Token string `json:"token"`
			}{Token: token}, rw)
	}
}

func (uh *UserHandler) FollowUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userInfo := r.Context().Value(keyValue{}).(userInformation)
		if err := uh.repo.FollowUser(userInfo.Username, getUserName(r)); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMessage{err.Error()}, rw)
			return
		}
		rw.WriteHeader(http.StatusNoContent)

	}
}

func (uh *UserHandler) UnFollowUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userInfo := r.Context().Value(keyValue{}).(userInformation)
		if err := uh.repo.UnFollowUser(userInfo.Username, getUserName(r)); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMessage{err.Error()}, rw)
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}

func (uh *UserHandler) GetLoggedInUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		client := r.Context().Value(keyValue{}).(userInformation)
		data.ToJSON(&client, rw)
	}
}

func (uh *UserHandler) HealthCheck() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		data.ToJSON(&generalMessage{"Good to go"}, rw)
	}
}

func getUserName(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["username"]
}
