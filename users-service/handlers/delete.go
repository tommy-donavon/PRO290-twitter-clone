package handlers

import (
	"net/http"

	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/users-service/data"
)

func (uh *UserHandler) DeleteUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userInformation := r.Context().Value(keyValue{}).(userInformation)
		userRecord, err := uh.repo.GetUser(userInformation.Username)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&generalMessage{"User not found"}, rw)
			return
		}
		if err := uh.repo.DeleteUser(userRecord.Username); err != nil {
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
