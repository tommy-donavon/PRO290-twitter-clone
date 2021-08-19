package handlers

import (
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/jasonlvhit/gocron"
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
	userInformation struct {
		Username string `json:"username"`
		UserType int    `json:"user_type"`
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

func (nh *NotificationHandler) NotificationConnection() *socketio.Server {
	nh.log.Println("GET SOCKET")
	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
			&websocket.Transport{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		},
	})
	getNotes := gocron.NewScheduler()
	sc := getNotes.Start()
	userInfo := userInformation{}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		nh.log.Println("connected:", s.ID())
		s.Emit("authorization")
		return nil
	})

	server.OnEvent("/", "request", func(s socketio.Conn, msg string) {
		resp, err := nh.SendNewRequest("users-service", "GET", "", map[string]string{"Authorization": msg})
		if err != nil {
			nh.log.Println(err)
		}

		if err := data.FromJSON(&userInfo, resp.Body); err != nil {
			nh.log.Println(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			nh.log.Println("Unauthorized user")
			s.Close()
			return
		}
		getNotes.Every(5).Seconds().Do(func() {
			nh.log.Println("I'm being run")
			nc, err := nh.repo.RetrieveNotifications(userInfo.Username)
			if err != nil {
				nh.log.Println(err)
			}
			notes := nc.Notification
			for _, v := range notes {
				s.Emit("notification", v)
				notes = notes[1:]
			}
			nc.Notification = notes
			err = nh.repo.SaveNotifications(nc)
			if err != nil {
				nh.log.Println(err)
			}
		})
		<-sc
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("[ERROR]:", e)

		getNotes.Clear()
		close(sc)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("[CLOSED]:", reason)
		getNotes.Clear()
		close(sc)

	})
	return server
}

func (nh *NotificationHandler) HealthCheck() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		data.ToJSON(&struct {
			Message string `json:"message"`
		}{"service good to go"}, rw)
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
