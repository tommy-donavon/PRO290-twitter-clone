package handlers

import (
	"context"
	"net/http"

	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/post-service/data"
)

func (ph *PostHandler) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			rw.WriteHeader(http.StatusForbidden)
			data.ToJSON(&generalMesage{"No token provided"}, rw)
			return
		}
		serr, err := ph.reg.LookUpService("users-service")
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMesage{err.Error()}, rw)
			return
		}
		req, err := http.NewRequest("GET", serr.GetHTTP(), nil)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMesage{err.Error()}, rw)
			return
		}
		req.Header.Add("Authorization", token)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMesage{err.Error()}, rw)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			rw.WriteHeader(http.StatusUnauthorized)
			data.ToJSON(&generalMesage{"You are not authorized to make this request"}, rw)
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
			data.ToJSON(&generalMesage{err.Error()}, rw)
			return
		}
		if err := post.Validate(); err != nil {
			ph.log.Println("[ERROR] validating item", err)
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&generalMesage{err.Error()}, rw)
			return
		}
		ctx := context.WithValue(r.Context(), keyValue{}, post)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}
