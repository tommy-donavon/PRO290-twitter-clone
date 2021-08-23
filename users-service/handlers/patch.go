package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/users-service/data"
)

func (uh *UserHandler) UpdateUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		uh.log.Println("PATCH USER")
		userInformation := r.Context().Value(keyValue{}).(userInformation)
		requestBody := map[string]string{}
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			uh.log.Println("[ERROR] unable to parse request body to map", err)
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&generalMessage{"Unable to process request body"}, rw)
			return
		}
		if err := uh.repo.UpdateUser(userInformation.Username, requestBody); err != nil {
			uh.log.Println("[ERROR] unable to update user", err)
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMessage{"unable to save user information"}, rw)
			return
		}
		if _, ok := requestBody["profile_uri"]; ok {
			ser, err := uh.reg.LookUpService("post-service")
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				data.ToJSON(&generalMessage{"unable to contact post service"}, rw)
				return
			}
			req, _ := http.NewRequest("PATCH", ser.GetHTTP()+"uri", nil)
			req.Header.Set("Authorization", r.Header.Get("Authorization"))
			client := &http.Client{}
			resp, _ := client.Do(req)
			if resp.StatusCode != http.StatusNoContent {
				rw.WriteHeader(http.StatusInternalServerError)
				data.ToJSON(&generalMessage{"o"}, rw)
				return
			}

		}
		rw.WriteHeader(http.StatusNoContent)
	}
}
