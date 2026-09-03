package websocket

import (
	"time"

	gws "github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	hub  *Hub
	conn *gws.Conn

	send chan []byte
}

func NewClient(
	hub *Hub,
	conn *gws.Conn,
) *Client {

	return &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
}

func (c *Client) readPump() {

	defer func() {
		c.hub.unregister <- c

		c.conn.Close()
	}()

	c.conn.SetReadDeadline(
		time.Now().Add(pongWait),
	)

	c.conn.SetPongHandler(func(string) error {

		c.conn.SetReadDeadline(
			time.Now().Add(pongWait),
		)

		return nil
	})

	for {

		_, _, err := c.conn.ReadMessage()

		if err != nil {
			break
		}
	}
}

func (c *Client) writePump() {

	ticker := time.NewTicker(
		pingPeriod,
	)

	defer func() {

		ticker.Stop()

		c.conn.Close()
	}()

	for {

		select {

		case message, ok := <-c.send:

			c.conn.SetWriteDeadline(
				time.Now().Add(writeWait),
			)

			if !ok {
				c.conn.WriteMessage(
					gws.CloseMessage,
					[]byte{},
				)

				return
			}

			err := c.conn.WriteMessage(
				gws.TextMessage,
				message,
			)

			if err != nil {
				return
			}

		case <-ticker.C:

			c.conn.SetWriteDeadline(
				time.Now().Add(writeWait),
			)

			if err := c.conn.WriteMessage(
				gws.PingMessage,
				nil,
			); err != nil {
				return
			}
		}
	}
}