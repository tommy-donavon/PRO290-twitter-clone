package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/post-service/handlers"
)

func SetUpRoutes(sm *mux.Router, postHandler *handlers.PostHandler) {
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", postHandler.CreatePost())
	postRouter.Use(postHandler.Auth)
	postRouter.Use(postHandler.MiddleWareValidatePost)
}
