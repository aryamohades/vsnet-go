package main

import (
	"github.com/aryamohades/sock/sock/epoll"
	"github.com/aryamohades/sock/stats"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"log"
	"net/http"
)

var ec *epoll.EventsCollector

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)

	if err != nil {
		log.Printf("upgrade error: %v", err)
		return
	}

	stats.AddConnection()

	if err := ec.Add(conn); err != nil {
		log.Printf("add connection error: %v", err)
		stats.RemoveConnection()

		if err := conn.Close(); err != nil {
			log.Printf("socket close error: %v", err)
		}
	}
}

func Start() {
	for {
		connections, err := ec.Wait()

		if err != nil {
			log.Printf("epoll wait error: %v", err)
			continue
		}

		for _, conn := range connections {
			if conn == nil {
				break
			}

			if _, _, err := wsutil.ReadClientData(conn); err != nil {
				if err := ec.Remove(conn); err != nil {
					log.Printf("remove error: %v", err)
				}

				stats.RemoveConnection()

				if err := conn.Close(); err != nil {
					log.Printf("socket close error: %v", err)
				}
			} else {
				// This is commented out since in demo usage, stdout is showing messages sent from > 1M connections at very high rate
				//log.Printf("msg: %s", string(msg))
			}
		}
	}
}

func main() {
	go stats.Start(":6000")

	var err error
	ec, err = epoll.New()

	if err != nil {
		log.Fatalf("new epoll error: %v", err)
	}

	go Start()

	mux := http.NewServeMux()
	mux.HandleFunc("/", wsHandler)

	srv := &http.Server{
		Handler: mux,
		Addr:    ":3333",
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("serve error: %v", err)
	}
}
