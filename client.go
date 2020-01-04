package main

import (
	"flag"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

var (
	ip          = flag.String("ip", "localhost", "ws service ip")
	port        = flag.String("port", "3333", "ws service port")
	connections = flag.Int("conn", 1, "number of websocket connections to open")
	ramp        = flag.Int("ramp", 0, "amount of seconds it takes to create all connections")
)

func main() {
	flag.Parse()

	u := url.URL{Scheme: "ws", Host: *ip + ":" + *port, Path: "/"}

	var conns []*websocket.Conn

	if *ramp <= 0 {
		for i := 0; i < *connections; i++ {
			c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

			if err != nil {
				log.Printf("connection %d failed: %v", i+1, err)
				break
			}

			conns = append(conns, c)
		}
	} else {
		for i := 0; i < *ramp; i++ {
			for j := 0; j < *connections / *ramp; j++ {
				c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

				if err != nil {
					log.Printf("connection %d failed: %v", i+1, err)
					break
				}

				conns = append(conns, c)

				defer func() {
					c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second))
					time.Sleep(time.Second)
					c.Close()
				}()
			}

			time.Sleep(1 * time.Second)
		}
	}

	log.Printf("initialized %d connections", len(conns))

	for {
		for i := 0; i < len(conns); i++ {
			conn := conns[i]

			if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(5*time.Second)); err != nil {
				log.Printf("pong error: %v", err)
			}

			if err := conn.WriteMessage(websocket.TextMessage, []byte(time.Now().String())); err != nil {
				log.Printf("write error: %v", err)
			}

			time.Sleep(1 * time.Second)
		}
	}
}
