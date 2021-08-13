package handlers

import (
	"net/http"

	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/users-service/data"
)

func (uh *UserHandler) GetLoggedInUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		client := r.Context().Value(keyValue{}).(userInformation)
		data.ToJSON(&client, rw)
	}
}
func (uh *UserHandler) GetFollowingList() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		username := getUserName(r)
		users, err := uh.repo.GetFollowingList(username)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMessage{err.Error()}, rw)
			return
		}
		data.ToJSON(users, rw)
	}
}
func (uh *UserHandler) GetFollowersList() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		username := getUserName(r)
		users, err := uh.repo.GetFollowersList(username)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMessage{err.Error()}, rw)
			return
		}
		data.ToJSON(users, rw)
	}
}

func (uh *UserHandler) HealthCheck() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		data.ToJSON(&generalMessage{"Good to go"}, rw)
	}
}
