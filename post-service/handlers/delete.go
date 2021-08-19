package handlers

import (
	"net/http"

	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/post-service/data"
)

func (ph *PostHandler) DeletePost() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ph.log.Println("DELETE POST")
		post := ph.repo.GetPost(uint(getPostId(r)))
		userInfo, err := ph.getUserInformation(r)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMesage{"Failed to connect to user api"}, rw)
			return
		}
		if post.Author != userInfo.Username {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&generalMesage{"You are not authorized to delete this post"}, rw)
			return
		}
		if err := ph.repo.DeletePost(post.ID); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMesage{"Unable to delete post with the provied id"}, rw)
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}
