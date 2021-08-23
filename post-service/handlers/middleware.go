package handlers

import (
	"context"
	"net/http"

	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/post-service/data"
)

func (ph *PostHandler) GlobalContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func (ph *PostHandler) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		resp, err := ph.sendNewRequest("users-service", "GET", "", map[string]string{
			"Authorization": r.Header.Get("Authorization"),
		})
		if err != nil {

			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMessage{err.Error()}, rw)
			return
		}

		if resp.StatusCode != http.StatusOK {
			rw.WriteHeader(http.StatusUnauthorized)
			data.ToJSON(&generalMessage{"You are not authorized to make this request"}, rw)
			return
		}
		next.ServeHTTP(rw, r)
	})
}

func (ph *PostHandler) MiddleWareValidatePost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		post := data.Post{}
		if err := data.FromJSON(&post, r.Body); err != nil {
			ph.log.Println("[ERROR] deserializing item", err)
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&generalMessage{err.Error()}, rw)
			return
		}
		if err := post.Validate(); err != nil {
			ph.log.Println("[ERROR] validating item", err)
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&generalMessage{err.Error()}, rw)
			return
		}
		ctx := context.WithValue(r.Context(), keyValue{}, post)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}
