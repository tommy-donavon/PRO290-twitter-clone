package handlers

import (
	"net/http"

	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/users-service/data"
)

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
