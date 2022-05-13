package main

import (
	"fmt"
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

var socketsAndUsers = make(map[int][]*websocket.Conn)

func initSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: websocket")

	// var userData UserSock

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		httpError(&w, 500, "error connecting to socket")
		return
	}

	// userData.Socket = c

	// expecting handshake to know what is the user id

	_, message, err := c.ReadMessage()

	if err != nil {
		httpError(&w, 500, "error reading message")
		return
	}

	Debugln(string(message[:]))

	userId, err := strconv.Atoi(string(message[:]))
	if err != nil {
		httpError(&w, 300, "id is not an int")
		return
	}

	// Debugln(userId)

	// socketsAndUsers = append(socketsAndUsers, userData)
	// 	tmparr := append(socketsAndUsers[userId], c)

	if socketsAndUsers[userId] == nil {
		socketsAndUsers[userId] = []*websocket.Conn{c}
	} else {
		socketsAndUsers[userId] = append(socketsAndUsers[userId], c)
	}

	// socketsAndUsers[userId] = "ciao"

}

func socketSendNotification(user int, str string) {

	sockArr, found := socketsAndUsers[user]

	if !found {
		return
	}

	for _, singleSocket := range sockArr {
		err := singleSocket.WriteMessage(websocket.TextMessage, []byte(str))
		if err != nil {
			Debugln("error in web socket: " + err.Error())
		}
	}

	return
}
