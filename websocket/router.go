package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func SetHandler(rt *mux.Router) {
	rt.HandleFunc("/msg", wsMsgHandler)
	rt.HandleFunc("/mem", wsMemHandler)
}

func wsMsgHandler(w http.ResponseWriter, r *http.Request) {
	chatId := r.URL.Query().Get("chatId")
	//userName := r.URL.Query().Get("userName")
	hub, ok := MsgHubs[chatId]
	if !ok {
		hub = newHub()
		go hub.Run()
		MsgHubs[chatId] = hub
	}
	connToHub(hub, "", w, r)
}

func wsMemHandler(w http.ResponseWriter, r *http.Request) {
	chatId := r.URL.Query().Get("chatId")
	userName := r.URL.Query().Get("userName")
	hub, ok := MemHubs[chatId]
	if !ok {
		hub = newMemHub()
		go hub.Run()
		MemHubs[chatId] = hub
	}
	connToHub(hub, userName, w, r)
}

func connToHub(hub IHub, userName string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrader connect: ", err)
		return
	}
	client := &Client{
		name: userName,
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	go client.writePump()
	go client.readPump()

	client.hub.Register() <- client
}
