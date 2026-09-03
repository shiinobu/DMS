package websocket

import (
	"net/http"

	gws "github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
)

var upgrader = gws.Upgrader{

	ReadBufferSize:  1024,

	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool {

		return true
	},
}

func Handler(hub *Hub) gin.HandlerFunc {

	return func(c *gin.Context) {

		conn, err := upgrader.Upgrade(
			c.Writer,
			c.Request,
			nil,
		)

		if err != nil {
			return
		}

		client := NewClient(
			hub,
			conn,
		)

		hub.register <- client

		go client.writePump()

		go client.readPump()
	}
}