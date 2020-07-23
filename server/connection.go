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
		c.Conn.WriteMessage(websocket.BinaryMessage, data)
	} else {
		c.Conn.WriteMessage(websocket.TextMessage, data)
	}
	return nil
}

func (c *WSConnection) WriteCloseMessage(code int, text string) error {
	if code == 0 {
		code = websocket.CloseNormalClosure
	}
	return c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, text))
}

func (c *WSConnection) DispatchMessages(onRead func(c *WSConnection, isBinary bool, data []byte)) error {
	for {
		ft, data, err := c.ReadMessage()
		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				fmt.Printf("connection [%s] closed passively: %v\n", c.GetID(), err)
				if c.Conn != nil {
					c.Conn.Close()
				}
				return err
			} else if err != nil {
				fmt.Printf("error occurred when reading message: %v, close the connection [%s]", err, c.GetID())
				c.Close(websocket.CloseNormalClosure, err.Error())
				return err
			}
		} else if ft == websocket.BinaryMessage {
			onRead(c, true, data)
		} else if ft == websocket.TextMessage {
			onRead(c, false, data)
		}
	}
	return nil
}

func (c *WSConnection) Close(code int, text string) {
	if c.Conn != nil {
		_ = c.CloseHandler()(websocket.CloseNormalClosure, text)
		c.Conn.Close()
	}
}
