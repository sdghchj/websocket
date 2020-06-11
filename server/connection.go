package server

import "github.com/gorilla/websocket"

type WSConnection struct {
	*websocket.Conn
	done chan bool
}

func NewWSConnection(conn *websocket.Conn) *WSConnection {
	return &WSConnection{
		Conn: conn,
		done: make(chan bool),
	}
}

func (c *WSConnection) Done() chan bool {
	return c.done
}

func (c *WSConnection) Close() {
	if c.done != nil {
		close(c.done)
		c.done = nil
	}
	if c.Conn != nil {
		c.Conn.Close()
		c.Conn = nil
	}
}
