package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

type Notification struct {
	UserID string `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func sendNotify(conn *nats.Conn, subject string, msg string) {
	if err := conn.Publish(subject, []byte(msg)); err != nil {
		log.Fatal(err)
	}
}

func main() {
	staticToken := os.Getenv("X_TOKEN")
	natsHosts := os.Getenv("NATS_HOSTS")
	natsSubject := os.Getenv("NATS_SUBJECT")

	if staticToken == "" || natsHosts == "" || natsSubject == "" {
		panic("Specify env vars X_TOKEN, NATS_HOSTS, NATS_SUBJECT")
	}

	// nats connection
	nc, err := nats.Connect(natsHosts,
		nats.Name("Testers tool"),
		nats.Timeout(10*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// server
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {})

	mux.HandleFunc("/send_notify", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-TOKEN")

		if token != staticToken {
			http.Error(w, "Access Denied", 403)
			return
		}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Bad json", 403)
		}
		s := string(b)
		sendNotify(nc, natsSubject, s)
	})
	log.Fatal(http.ListenAndServe(":8080", mux))
}
