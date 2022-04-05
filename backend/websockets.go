package main

import (
	"net/http"

	"github.com/gorilla/websocket"

	_ "github.com/go-sql-driver/mysql"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func initSocket(w http.ResponseWriter, r *http.Request) {
	Debugln("endpoint hit: websocket")

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		httpError(&w, 500, "error connecting to socket")
		return
	}

	for {

		messageType, message, err := c.ReadMessage()
		if err != nil {
			httpError(&w, 500, "error reading message")
			return
		}

		Debugln(message)

		if err = c.WriteMessage(messageType, message); err != nil {
			httpError(&w, 500, "error sending message")
			return
		}
	}

}
