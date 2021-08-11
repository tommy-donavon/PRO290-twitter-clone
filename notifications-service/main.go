package main

import (
	"log"
	"net/http"
	"os"

	socketio "github.com/googollee/go-socket.io"
	"github.com/yhung-mea7/PRO290-twitter-clone/blob/main/notifications-service/data"
	"github.com/yhung-mea7/PRO290-twitter-clone/blob/main/notifications-service/register"
)

func main() {
	server := socketio.NewServer(nil)

	consulClient := register.NewConsulClient("notifications-service")
	consulClient.RegisterService()
	defer consulClient.DeregisterService()

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		return nil
	})

	// server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
	// 	log.Println("notice:", msg)
	// 	s.Emit("reply", "have "+msg)
	// })

	// server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
	// 	s.SetContext(msg)
	// 	return "recv " + msg
	// })

	// server.OnEvent("/", "bye", func(s socketio.Conn) string {
	// 	last := s.Context().(string)
	// 	s.Emit("bye", last)
	// 	s.Close()
	// 	return last
	// })

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed", reason)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer server.Close()

	http.Handle("/", server)
	http.Handle("/healthcheck", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		data.ToJSON(&struct{ message string }{"service good to go"}, rw)
	}))

	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), nil))
}
