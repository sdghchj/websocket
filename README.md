# websocket
websocket wrapper of  https://github.com/gorilla/websocket

**A websocker server handler**

Example:

```
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sdghchj/websocket/server"
)

func main() {
	handler := server.NewWSServerHandler(
		func(conn *server.WSConnection) {
			fmt.Println("a connection opened")
		},
		func(conn *server.WSConnection) {
			fmt.Println("the peer closed the websocket connection")
		},
		func(conn *server.WSConnection, isBinary bool, data []byte) {
			if isBinary {
				fmt.Println("read:", data)
			} else {
				fmt.Println("read:", string(data))
			}
			conn.WriteMessage(false, []byte("hello, received your data"))
		})
	err := http.ListenAndServe("0.0.0.0:80", handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
```
