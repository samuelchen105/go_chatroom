package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var Hubs = make(map[string]*Hub)

//MemHubs = make(map[string]*MemHub)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func SetHandler(rt *mux.Router) {
	rt.HandleFunc("/", wsHandler)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	chatId := r.URL.Query().Get("chatId")
	userName := r.URL.Query().Get("userName")
	hub, ok := Hubs[chatId]
	if !ok {
		hub = newHub()
		go hub.Run()
		Hubs[chatId] = hub
	}
	connToHub(hub, userName, w, r)
}

func connToHub(hub *Hub, userName string, w http.ResponseWriter, r *http.Request) {
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
