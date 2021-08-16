package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yhung-mea7/PRO290-twitter-clone/blob/main/notifications-service/handlers"
)

func SetUpRoutes(sm *mux.Router, nh *handlers.NotificationHandler) {
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/healthcheck", nh.HealthCheck())

}
