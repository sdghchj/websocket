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
	getID          func(conn *http.Request) string
}

func NewWSServerHandler(onRead func(c *WSConnection, isBinary bool, data []byte)) *WSServerHandler {
	return &WSServerHandler{
		onRead: onRead,
	}
}

func (handler *WSServerHandler) SetConnectHandler(onConnect func(conn *WSConnection)) *WSServerHandler {
	handler.onConnect = onConnect
	return handler
}

func (handler *WSServerHandler) SetCloseHandler(onClose func(conn *WSConnection)) *WSServerHandler {
	handler.onClose = onClose
	return handler
}

func (handler *WSServerHandler) SetIDGetter(getID func(conn *http.Request) string) *WSServerHandler {
	handler.getID = getID
	return handler
}

func (handler *WSServerHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if handler.onRead == nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	var id string
	if handler.getID != nil {
		id = handler.getID(request)
	}

	conn, err := handler.Upgrader.Upgrade(response, request, handler.responseHeader)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}

	wsc := NewWSConnection(conn, id)

	if handler.onClose != nil {
		defaultCloser := conn.CloseHandler()
		conn.SetCloseHandler(func(code int, text string) error {
			defaultCloser(code, text)
			if handler.onClose != nil {
				handler.onClose(wsc)
			}
			return nil
		})
	}

	if handler.onConnect != nil {
		handler.onConnect(wsc)
	}

	go wsc.dispatchMessages(handler.onRead)
}
