package server

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type WSConnection struct {
	*websocket.Conn
	id string
}

func NewWSConnection(conn *websocket.Conn, id string) *WSConnection {
	return &WSConnection{
		Conn: conn,
		id:   id,
	}
}

func (c *WSConnection) GetID() string {
	return c.id
}

func (c *WSConnection) WriteMessage(isBinary bool, data []byte) error {
	if isBinary {
		return c.Conn.WriteMessage(websocket.BinaryMessage, data)
	}
	return c.Conn.WriteMessage(websocket.TextMessage, data)
}

func (c *WSConnection) dispatchMessages(onRead func(c *WSConnection, isBinary bool, data []byte)) error {
	defer c.Conn.Close()
	for {
		ft, data, err := c.Conn.ReadMessage()
		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				//OnClose has run in ReadMessage
				fmt.Printf("connection [%s] closed passively: %v\n", c.GetID(), err)
				return err
			} else if err != nil {
				fmt.Printf("error occurred when reading message: %v, close the connection [%s]", err, c.GetID())
				c.Close(websocket.CloseAbnormalClosure, err.Error())
			}
		} else if ft == websocket.BinaryMessage {
			onRead(c, true, data)
		} else if ft == websocket.TextMessage {
			onRead(c, false, data)
		}
	}
}

func (c *WSConnection) Close(code int, text string) error {
	return c.CloseHandler()(code, text) //including OnClose
}
