package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/post-service/handlers"
)

func SetUpRoutes(sm *mux.Router, postHandler *handlers.PostHandler) {
	sm.Use(postHandler.GlobalContentTypeMiddleware)
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", postHandler.CreatePost())
	postRouter.HandleFunc("/{id:[0-9]+}", postHandler.CreatePost())
	postRouter.Use(postHandler.Auth)
	postRouter.Use(postHandler.MiddleWareValidatePost)

	feedRouter := sm.Methods(http.MethodGet).Subrouter()
	feedRouter.HandleFunc("/feed", postHandler.GetFeed())
	feedRouter.Use(postHandler.Auth)

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/healthcheck", postHandler.HealthCheck())
	getRouter.HandleFunc("/{id:[0-9]+}", postHandler.GetPost())
	getRouter.HandleFunc("/", postHandler.GetAllPosts())

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}", postHandler.DeletePost())
	deleteRouter.Use(postHandler.Auth)

	patchRouter := sm.Methods(http.MethodPatch).Subrouter()
	patchRouter.HandleFunc("/{id:[0-9]+}/like", postHandler.LikePost())
	patchRouter.HandleFunc("/{id:[0-9]+}/unlike", postHandler.UnlikePost())
	patchRouter.Use(postHandler.Auth)
}
