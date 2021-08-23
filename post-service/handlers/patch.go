package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/post-service/amqp"
	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/post-service/data"
)

func (ph *PostHandler) LikePost() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userInfo, err := ph.getUserInformation(r)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMessage{"unable to connect to user service"}, rw)
			return
		}
		post := ph.repo.GetPost(uint(getPostId(r)))
		if post.Author == "" {
			ph.log.Println("[ERROR] No post found")
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&generalMessage{"No post found"}, rw)
			return
		}

		if err := ph.repo.LikePost(post.ID); err != nil {
			ph.log.Println("[ERROR] unable to like post")
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMessage{"Unable to like post"}, rw)
			return
		}

		rw.WriteHeader(http.StatusNoContent)
		ph.messenger.SubmitToMessageBroker(&amqp.Message{
			Username: post.Author,
			Message:  fmt.Sprintf("%s liked your post!", userInfo.Username),
		})

	}
}

func (ph *PostHandler) UnlikePost() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		post := ph.repo.GetPost(uint(getPostId(r)))
		if post.Author == "" {
			ph.log.Println("[ERROR] No post found")
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&generalMessage{"No post found"}, rw)
			return
		}
		if err := ph.repo.UnlikePost(post.ID); err != nil {
			ph.log.Println("[ERROR] unable to unlike post", err)
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMessage{"Unable to unlike post"}, rw)
			return
		}

	}
}

func (ph *PostHandler) UpdatePost() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userInfo, err := ph.getUserInformation(r)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMessage{"unable to reach user service"}, rw)
			return
		}
		requestBody := map[string]string{}
		err = json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&generalMessage{"Unable to process request body"}, rw)
			return
		}
		post := ph.repo.GetPost(uint(getPostId(r)))
		if post.Author != userInfo.Username {
			rw.WriteHeader(http.StatusForbidden)
			data.ToJSON(&generalMessage{"you are not allowed to edit this post"}, rw)
			return
		}
		if err := ph.repo.UpdatePost(post.ID, requestBody); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMessage{"unable to save post information"}, rw)
			return
		}
		rw.WriteHeader(http.StatusNoContent)

	}
}

func (ph *PostHandler) UpdateAllAuthorUri() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userInfo, err := ph.getUserInformation(r)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMessage{"unable to reach user service"}, rw)
			return
		}
		requestBody := map[string]string{}
		err = json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&generalMessage{"Unable to process request body"}, rw)
			return
		}
		if err := ph.repo.UpdateAllAuthorUri(userInfo.Username, userInfo.ProfileUri); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMessage{"unable to update post uri"}, rw)
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}
