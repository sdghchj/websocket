package main

import (
	"flag"
	"github.com/sdghchj/websocket/client"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var addr = flag.String("addr", "ws://localhost:80", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, _, err := client.NewWSConnection(*addr, http.Header{"client-id": []string{"1"}}, func(c *client.WSConnection, isBinary bool, data []byte) {
		if isBinary {
			log.Printf("recv: %v", data)
		} else {
			log.Printf("recv: %s", string(data))
		}
	})
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			err := c.WriteMessage(false, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			return
		}
	}
}
