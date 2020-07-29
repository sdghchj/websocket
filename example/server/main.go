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
