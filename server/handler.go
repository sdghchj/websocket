package server

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type WSServerHandler struct {
	websocket.Upgrader
	responseHeader http.Header
	onConnect      func(conn *WSConnection)
	onClose        func(conn *WSConnection)
	onRead         func(c *WSConnection, isBinary bool, data []byte)
}

func NewWSServerHandler(onConnect, onclose func(conn *WSConnection), onRead func(c *WSConnection, isBinary bool, data []byte)) *WSServerHandler {
	return &WSServerHandler{
		onConnect: onConnect,
		onClose:   onclose,
		onRead:    onRead,
	}
}

func (handler *WSServerHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if handler.onConnect == nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	conn, err := handler.Upgrader.Upgrade(response, request, handler.responseHeader)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}

	wsc := NewWSConnection(conn)

	if handler.onClose != nil {
		defaultCloser := conn.CloseHandler()
		conn.SetCloseHandler(func(code int, text string) error {
			defaultCloser(code, text)
			handler.onClose(wsc)
			return nil
		})
	}

	handler.onConnect(wsc)

	go wsc.DispatchMessages(handler.onRead)
}
