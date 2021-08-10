package handlers

import (
	"net/http"

	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/post-service/data"
)

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

func (ph *PostHandler) GetAllPosts() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ph.log.Println("GET ALL POSTS")
		rw.Header().Set("Content-type", "application/json")
		data.ToJSON(ph.repo.GetAllPosts(), rw)
	}
}

func (ph *PostHandler) HealthCheck() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		data.ToJSON(&generalMesage{"Good to go"}, rw)
	}
}
