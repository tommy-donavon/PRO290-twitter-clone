package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/post-service/amqp"
	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/post-service/data"
	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/post-service/register"
)

type (
	PostHandler struct {
		repo      *data.PostRepo
		log       *log.Logger
		reg       *register.ConsulClient
		messenger *amqp.Messenger
	}
	generalMesage struct {
		Message interface{} `json:"message"`
	}
	userInformation struct {
		Username string `json:"username"`
		UserType int    `json:"user_type"`
	}
	keyValue struct{}
)

func NewPostHandler(repo *data.PostRepo, log *log.Logger, reg *register.ConsulClient, messenger *amqp.Messenger) *PostHandler {
	return &PostHandler{repo, log, reg, messenger}
}

func (ph *PostHandler) CreatePost() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		post := r.Context().Value(keyValue{}).(data.Post)
		userInfo, err := ph.getUserInformation(r)
		if err != nil {
			ph.log.Println("[ERROR] Unable to establish connection to User service", err)
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMesage{"Unable to establish connection to external service"}, rw)
			return
		}
		post.Author = userInfo.Username
		id := getPostId(r)
		if err := ph.repo.CreatePost(&post, uint(id)); err != nil {
			ph.log.Println("[ERROR] saving post in database", err)
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMesage{err.Error()}, rw)
			return
		}
		if id != 0 {
			if err := ph.messenger.SubmitToMessageBroker(&amqp.Message{
				Username: ph.repo.GetPost(uint(id)).Author,
				Message:  fmt.Sprintf("%s commented %s on your post", userInfo.Username, post.PostBody),
			}); err != nil {
				ph.log.Println("[ERROR] Sending message to queue failed", err)
			}
		}
		rw.WriteHeader(http.StatusNoContent)

	}
}

func (ph *PostHandler) getUserInformation(r *http.Request) (*userInformation, error) {
	resp, err := ph.sendNewRequest("users-service", "GET", "", map[string]string{"Authorization": r.Header.Get("Authorization")})
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}
	defer resp.Body.Close()
	userInfo := &userInformation{}
	if err := data.FromJSON(&userInfo, resp.Body); err != nil {
		return nil, err
	}
	return userInfo, nil
}

func (ph *PostHandler) sendNewRequest(serviceName, methodType, endpoint string, headerOptions map[string]string) (*http.Response, error) {
	ser, err := ph.reg.LookUpService(serviceName)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(methodType, ser.GetHTTP()+endpoint, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range headerOptions {
		req.Header.Set(key, value)
	}
	client := &http.Client{}
	return client.Do(req)
}

func getPostId(r *http.Request) int {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return 0
	}
	return id
}
