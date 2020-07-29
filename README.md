# websocket
websocket wrapper of  https://github.com/gorilla/websocket

**A websocker server handler**

Example:

Server:
```golang
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sdghchj/websocket/server"
)

func main() {
	handler := server.NewWSServerHandler(
		func(conn *server.WSConnection, isBinary bool, data []byte) {
			if isBinary {
				fmt.Println("read:", data)
			} else {
				fmt.Println("read:", string(data))
			}
			conn.WriteMessage(false, []byte("hello, received your data"))
		}).SetConnectHandler(func(conn *server.WSConnection) {
		fmt.Printf("a connection opened: %s\n", conn.GetID())
	}).SetCloseHandler(func(conn *server.WSConnection) {
		fmt.Printf("the peer closed the websocket connection: %s\n", conn.GetID())
	}).SetIDGetter(func(request *http.Request) string {
		return request.Header.Get("client-id")
	})
	err := http.ListenAndServe("0.0.0.0:80", handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
```


Client:
```golang
import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

    "github.com/sdghchj/websocket/client"
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

```