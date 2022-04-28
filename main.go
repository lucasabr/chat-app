package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)


var (
	wsConn *websocket.Conn
	wsUpgrader = websocket.Upgrader { 
	ReadBufferSize: 1024,
	WriteBufferSize: 1024, 
}
)

type Message struct {
	Text string `json:"Text"`
	User string `json:"User"`
}

func wsHandler(w http.ResponseWriter,  r *http.Request){ 

	wsUpgrader.CheckOrigin = func(r *http.Request) bool{ 
		return true
	}

	var err error
	wsConn, err := wsUpgrader.Upgrade(w, r, nil)
	if err!= nil {
		fmt.Printf("error upgrading: %s\n", err.Error())
		return
	}

	defer wsConn.Close()

	for {
		var msg Message 
		
		err := wsConn.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("Error reading JSON %s\n", err.Error())
			break
		}

		fmt.Printf("Message Received: %s\n", msg.Text)
	}
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/socket", wsHandler)

	log.Fatal(http.ListenAndServe(":4500", router))
}