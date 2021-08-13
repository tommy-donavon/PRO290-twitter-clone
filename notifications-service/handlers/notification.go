package handlers

import (
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/yhung-mea7/PRO290-twitter-clone/blob/main/notifications-service/amqp"
	"github.com/yhung-mea7/PRO290-twitter-clone/blob/main/notifications-service/data"
	"github.com/yhung-mea7/PRO290-twitter-clone/blob/main/notifications-service/register"
)

type (
	NotificationHandler struct {
		log       *log.Logger
		repo      *data.NotificationRepo
		reg       *register.ConsulClient
		messenger *amqp.Messenger
		Server    *socketio.Server
	}
)

func NewNotificationHandler(repo *data.NotificationRepo, log *log.Logger, reg *register.ConsulClient, messenger *amqp.Messenger) *NotificationHandler {
	messenger.Consume()
	return &NotificationHandler{
		log:       log,
		repo:      repo,
		messenger: messenger,
		reg:       reg,
	}
}

// func (nh *NotificationHandler) NotificationConnection() {
// 	nh.log.Println("GET SOCKET")
// 	// nh.Server = socketio.NewServer(nil)

// 	server.OnConnect("/", func(s socketio.Conn) error {
// 		s.SetContext("")
// 		log.Println("connected:", s.ID())
// 		s.Emit("authorization")
// 		return nil
// 	})

// 	server.OnEvent("/", "request", func(s socketio.Conn, msg string) {
// 		resp, err := nh.sendNewRequest("users-service", "GET", "", map[string]string{"Authorization": msg})
// 		if err != nil {
// 			nh.log.Println(err)
// 		}
// 		userInfo := userInformation{}
// 		if err := data.FromJSON(&userInfo, resp.Body); err != nil {
// 			nh.log.Println(err)
// 		}
// 		resp.Body.Close()
// 		nc, err := nh.repo.RetrieveNotifications(userInfo.Username)
// 		if err != nil {
// 			nh.log.Println(err)
// 		}
// 		notes := nc.Notification
// 		for _, v := range notes {

// 			s.Emit("notification", v)
// 			// notes = notes[1:]
// 		}
// 		nc.Notification = notes
// 		err = nh.repo.SaveNotifications(nc)
// 		if err != nil {
// 			nh.log.Println(err)
// 		}

// 	})

// 	server.OnError("/", func(s socketio.Conn, e error) {
// 		log.Println("meet error:", e)
// 	})

// 	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
// 		log.Println("closed", reason)
// 	})
// 	go func() {
// 		if err := server.Serve(); err != nil {
// 			log.Fatalf("socketio listen error: %s\n", err)
// 		}
// 	}()

// }

func (nh *NotificationHandler) HealthCheck() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		data.ToJSON(&struct{ message string }{"service good to go"}, rw)
	}
}
func (ph *NotificationHandler) SendNewRequest(serviceName, methodType, endpoint string, headerOptions map[string]string) (*http.Response, error) {
	ser, err := ph.reg.LookUpService(serviceName)
	if err != nil {
		return nil, err
	}
	ph.log.Println(ser.GetHTTP() + endpoint)
	req, err := http.NewRequest(methodType, ser.GetHTTP()+endpoint, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range headerOptions {
		req.Header.Set(key, value)
	}
	client := &http.Client{}
	return client.Do(req)
}
