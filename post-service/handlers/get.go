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

func (ph *PostHandler) GetFeed() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ph.log.Println("GET FEED")
		rw.Header().Set("Content-type", "application/json")
		userInfo,err := ph.getUserInformation(r)
		if err != nil {
			ph.log.Println("[ERROR] Unable to establish connection to internal service", err)
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMesage{"Unable to establish connection to internal service"}, rw)
			return

		}
		resp, err := ph.sendNewRequest("users-service", "GET", userInfo.Username, map[string]string{"Authorization": r.Header.Get("Authorization")})
		if err != nil || resp.StatusCode != http.StatusOK {
			ph.log.Println("[ERROR] Unable to establish connection to internal service", err)
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMesage{"Unable to establish connection to internal service"}, rw)
			return
		}
		defer resp.Body.Close()
		followingList := []*data.FollowInformation{}
		if err := data.FromJSON(&followingList, resp.Body); err != nil {
			ph.log.Println("[ERROR] deserializing response body", err)
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&generalMesage{"Unable to retrieve user information"}, rw)
			return
		}
		data.ToJSON(ph.repo.GetFeed(userInfo.Username, followingList), rw)
	}
}

func (ph *PostHandler) HealthCheck() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		data.ToJSON(&generalMesage{"Good to go"}, rw)
	}
}
