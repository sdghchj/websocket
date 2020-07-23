package client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type WSConnection struct {
	*websocket.Conn
	done chan bool
}

func NewWSConnection(endpoint string, headers http.Header, onRead func(c *WSConnection, isBinary bool, data []byte)) (*WSConnection, *http.Response, error) {
	c, resp, err := websocket.DefaultDialer.Dial(endpoint, headers)
	if err != nil {
		return nil, resp, err
	}
	wc := &WSConnection{Conn: c, done: make(chan bool)}
	go wc.dispatchMessages(onRead)
	return wc, resp, err
}

func (c *WSConnection) dispatchMessages(onRead func(c *WSConnection, isBinary bool, data []byte)) error {
	defer func() {
		c.Conn.Close()
		close(c.done)
	}()
	for {
		ft, data, err := c.ReadMessage()
		if err != nil {
			if peerCloserErr, ok := err.(*websocket.CloseError); ok {
				fmt.Printf("connection onClose: %v\n", err)
				_ = c.WriteCloseMessage(peerCloserErr.Code, peerCloserErr.Text)
				return err
			} else if err != nil {
				fmt.Printf("error occurred when reading message, close the connection [%v]\n", err)
				_ = c.WriteCloseMessage(websocket.CloseAbnormalClosure, err.Error())
			}
		} else if ft == websocket.BinaryMessage {
			onRead(c, true, data)
		} else if ft == websocket.TextMessage {
			onRead(c, false, data)
		}
	}
}

func (c *WSConnection) Done() <-chan bool {
	return c.done
}

func (c *WSConnection) WriteMessage(isBinary bool, data []byte) error {
	if isBinary {
		return c.Conn.WriteMessage(websocket.BinaryMessage, data)
	}
	return c.Conn.WriteMessage(websocket.TextMessage, data)
}

func (c *WSConnection) WriteCloseMessage(code int, text string) error {
	if code == 0 {
		code = websocket.CloseNormalClosure
	}
	return c.Conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(code, text), time.Now().Add(time.Second*10))
}

func (c *WSConnection) Close() error {
	err := c.WriteCloseMessage(websocket.CloseNormalClosure, "")
	if err != nil && err != websocket.ErrCloseSent {
		return err
	}
	select {
	case <-time.NewTimer(time.Second * 10).C:
		return fmt.Errorf("timeout")
	case <-c.done:
		return nil
	}
}
