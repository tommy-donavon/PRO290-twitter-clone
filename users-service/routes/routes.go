package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/users-service/handlers"
)

func SetUpRoutes(sm *mux.Router, userHandler *handlers.UserHandler) {
	signUpRouter := sm.Methods(http.MethodPost).Subrouter()
	signUpRouter.HandleFunc("/sign-up", userHandler.CreateUser())
	signUpRouter.Use(userHandler.MiddlewareValidateUser)

	loginRouter := sm.Methods(http.MethodPost).Subrouter()
	loginRouter.HandleFunc("/", userHandler.LoginUser())
	loginRouter.Use(userHandler.MiddlewareValidateLogin)

	followRouter := sm.Methods(http.MethodPost).Subrouter()
	followRouter.HandleFunc("/{username}", userHandler.FollowUser())
	followRouter.Use(userHandler.Auth)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{username}", userHandler.UnFollowUser())
	deleteRouter.HandleFunc("/", userHandler.DeleteUser())
	deleteRouter.Use(userHandler.Auth)

	healthHandler := sm.Methods(http.MethodGet).Subrouter()
	healthHandler.HandleFunc("/healthcheck", userHandler.HealthCheck())

	getUserHandler := sm.Methods(http.MethodGet).Subrouter()
	getUserHandler.HandleFunc("/", userHandler.GetLoggedInUser())
	getUserHandler.Use(userHandler.Auth)

	getFollowers := sm.Methods(http.MethodGet).Subrouter()
	getFollowers.HandleFunc("/{username}", userHandler.GetFollowingList())

}
