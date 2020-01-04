package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aryamohades/sock/stats"
	"nhooyr.io/websocket"
)

var (
	port = flag.String("port", "3333", "ws service port")
)


func wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})

	if err != nil {
		log.Printf("upgrade error: %v", err)
		return
	}

	defer c.Close(websocket.StatusInternalError, "")
}

func main() {
	go stats.Start(":6000")

	flag.Parse()

	l, err := net.Listen("tcp", ":" + *port)

	if err != nil {
		log.Fatalf("listen error: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", wsHandler)

	srv := &http.Server{
		Handler: mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	if err := srv.Serve(l); err != http.ErrServerClosed {
		log.Fatalf("serve error: %v", err)
	}

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
}
