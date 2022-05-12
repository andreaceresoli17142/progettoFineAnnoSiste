package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"

	_ "github.com/go-sql-driver/mysql"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type UserSock struct {
	Id     int
	Socket *websocket.Conn
}

var socketsAndUsers []UserSock

func initSocket(w http.ResponseWriter, r *http.Request) {
	Debugln("endpoint hit: websocket")

	var userData UserSock

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		httpError(&w, 500, "error connecting to socket")
		return
	}

	userData.Socket = c

	// expecting handshake to know what is the user id

	_, message, err := c.ReadMessage()

	if err != nil {
		httpError(&w, 500, "error reading message")
		return
	}

	Debugln(string(message[:]))

	userData.Id, err = strconv.Atoi(string(message[:]))
	if err != nil {
		httpError(&w, 300, "id is not an int")
		return
	}

	socketsAndUsers = append(socketsAndUsers, userData)
}

func socketSendNotification(user int) error {
	for _, userSock := range socketsAndUsers {
		if userSock.Id == user {
			if err := userSock.Socket.WriteMessage(websocket.BinaryMessage, []byte{0}); err != nil {
				return err
			}
		}
	}
	return nil
}
