package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/post-service/data"
	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/post-service/register"
)

type (
	PostHandler struct {
		repo *data.PostRepo
		log  *log.Logger
		reg  *register.ConsulClient
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

func NewPostHandler(repo *data.PostRepo, log *log.Logger, reg *register.ConsulClient) *PostHandler {
	return &PostHandler{repo, log, reg}
}

func (ph *PostHandler) CreatePost() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		post := r.Context().Value(keyValue{}).(data.Post)
		resp, err := ph.sendNewRequest("users-service", "GET", map[string]string{"Authorization": r.Header.Get("Authorization")})
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
		post.Author = userInfo.Username
		id := getPostId(r)
		if err := ph.repo.CreatePost(&post, uint(id)); err != nil {
			ph.log.Println("[ERROR] saving post in database", err)
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMesage{err.Error()}, rw)
			return
		}
		rw.WriteHeader(http.StatusNoContent)

	}
}

func (ph *PostHandler) GetPost() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ph.log.Println("GET POST")
		rw.Header().Set("Content-type", "application/json")
		id := getPostId(r)
		post := ph.repo.GetPost(uint(id))
		if post.ID == 0 {
			rw.WriteHeader(http.StatusNotFound)
			data.ToJSON(&generalMesage{"Post not found"}, rw)
			return
		}
		data.ToJSON(post, rw)
	}
}

func (ph *PostHandler) HealthCheck() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		data.ToJSON(&generalMesage{"Good to go"}, rw)
	}
}

func (ph *PostHandler) sendNewRequest(serviceName, methodType string, headerOptions map[string]string) (*http.Response, error) {
	ser, err := ph.reg.LookUpService(serviceName)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(methodType, ser.GetHTTP(), nil)
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
		return 0 //do something better than panic
	}
	return id
}
