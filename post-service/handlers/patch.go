package handlers

import (
	"fmt"
	"net/http"

	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/post-service/amqp"
	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/post-service/data"
)

func (ph *PostHandler) LikePost() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		resp, err := ph.sendNewRequest("users-service", "GET", "", map[string]string{"Authorization": r.Header.Get("Authorization")})
		if err != nil || resp.StatusCode != http.StatusOK {
			ph.log.Println("[ERROR] Unable to establish connection to internal service", err)
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMesage{"Unable to establish connection to internal service"}, rw)
			return
		}
		defer resp.Body.Close()
		userInfo := &userInformation{}
		if err := data.FromJSON(&userInfo, resp.Body); err != nil {
			ph.log.Println("[ERROR] deserializing response body", err)
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&generalMesage{"Unable to retrieve user information"}, rw)
			return
		}
		post := ph.repo.GetPost(uint(getPostId(r)))
		if post.Author == "" {
			ph.log.Println("[ERROR] No post found", err)
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&generalMesage{"No post found"}, rw)
			return
		}

		if err := ph.repo.LikePost(post.ID); err != nil {
			ph.log.Println("[ERROR] unable to like post")
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMesage{"Unable to like post"}, rw)
			return
		}

		rw.WriteHeader(http.StatusNoContent)
		ph.messenger.SubmitToMessageBroker(&amqp.Message{
			Username: post.Author,
			Message:  fmt.Sprintf("%s liked your post!", userInfo.Username),
		})

	}
}
