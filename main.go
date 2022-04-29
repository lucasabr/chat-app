package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)


var (
	connToName = make(map[*websocket.Conn]string)
	userToConn = make(map[string]*websocket.Conn)
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

	params := mux.Vars(r)
	user := params["user"]
	connToName[wsConn] = user
	userToConn[user] = wsConn
	
	defer wsConn.Close()

	for {
		var msg Message 
		
		err := wsConn.ReadJSON(&msg)
		if err != nil {
			disconnect(wsConn)
			break
		}


		dispatchMessages(msg)
	}
}

func disconnect(ws *websocket.Conn) {
	user := connToName[ws]
	delete(connToName, ws)
	delete(userToConn, user)
}



func dispatchMessages(msg Message) {
	for key := range connToName {
		err := key.WriteJSON(msg)
		if err != nil {
			fmt.Printf("error sending msg: %s\n", err.Error())
			return
		}
	}

}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/socket/{user}", wsHandler)

	log.Fatal(http.ListenAndServe(":4500", router))
}