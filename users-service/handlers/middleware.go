package handlers

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/users-service/data"
)

func (uh *UserHandler) MiddlewareValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		User := data.User{}
		if err := data.FromJSON(&User, r.Body); err != nil {
			uh.log.Println(err)
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&generalMessage{err.Error()}, rw)
			return
		}
		if err := User.Validate(); err != nil {
			uh.log.Println(err)
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&generalMessage{formatValidationError(err.Error())}, rw)
			return
		}
		ctx := context.WithValue(r.Context(), keyValue{}, User)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})

}

func (uh *UserHandler) MiddlewareValidateLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		login := data.Login{}
		if err := data.FromJSON(&login, r.Body); err != nil {
			uh.log.Println("[ERROR] deserializing login", err)
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&generalMessage{err.Error()}, rw)
			return
		}
		if err := login.Validate(); err != nil {
			uh.log.Println("[ERROR] login validation failed")
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&generalMessage{formatValidationError(err.Error())}, rw)
			return
		}
		ctx := context.WithValue(r.Context(), keyValue{}, login)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})

}

func (uh *UserHandler) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		token := r.Header.Get("Authorization")
		if token == "" {
			rw.WriteHeader(http.StatusForbidden)
			data.ToJSON(&generalMessage{"No token provided"}, rw)
			return
		}
		jwToken := strings.Split(token, "Bearer ")
		if len(jwToken) == 2 {
			token = strings.TrimSpace(jwToken[1])
		} else {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&generalMessage{"Malformed Token"}, rw)
			return
		}
		claims, err := uh.jwt.CheckToken(token)
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			data.ToJSON(&generalMessage{"Unauthorized"}, rw)
			return
		}
		client := userInformation{claims.Username, claims.UserType}
		ctx := context.WithValue(r.Context(), keyValue{}, client)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)

	})
}

func formatValidationError(err string) map[string]string {
	messages := strings.Split(err, "\n")
	output := map[string]string{}
	rgx := regexp.MustCompile(`^Key: (.*) Error:(.*)$`)
	for _, line := range messages {
		m := rgx.FindStringSubmatch(line)
		output[m[1]] = m[2]

	}
	return output
}
