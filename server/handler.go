package server

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type WSServerHandler struct {
	websocket.Upgrader
	responseHeader http.Header
	onConnect      func(conn *WSConnection)
	onRead         func(conn *WSConnection, frameType int, data []byte)
	onClose        func(conn *WSConnection)
}

func NewWSServerHandler() *WSServerHandler {

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
		conn.SetCloseHandler(func(code int, text string) error {
			wsc.Close()
			message := websocket.FormatCloseMessage(code, "")
			wsc.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Minute))
			handler.onClose(wsc)
			return nil
		})
	}

	handler.onConnect(wsc)

	if handler.onRead != nil {
		go func() {
			for {
				select {
				case <-wsc.Done():
					return
				default:
					ft, data, err := wsc.ReadMessage()
					if err != nil {
						wsc.Close()
						return
					} else {
						handler.onRead(wsc, ft, data)
					}
				}
			}
		}()
	}
}
