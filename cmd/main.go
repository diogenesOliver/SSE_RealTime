package main

import (
	"fmt"
	"net/http"

	"github.com/diogenesOliver/sse_realTime/package/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	out := make(chan amqp.Delivery)
	rabbitmqChannel, err := rabbitmq.OpenChanel()
	if err != nil {
		panic(err)
	}
	go rabbitmq.Consume("msgs", rabbitmqChannel, out)

	http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		for m := range out {
			fmt.Fprintf(w, "event: message\n")
			fmt.Fprintf(w, "data: %s\n\n", m.Body)
			w.(http.Flusher).Flush()
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "html/inde.html")
	})
	http.ListenAndServe(":8080", nil)
}
