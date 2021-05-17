package common

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
	MsgHubs = make(map[string]*Hub)
	MemHubs = make(map[string]*MemHub)
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	name string
	hub  IHub
	conn *websocket.Conn
	send chan []byte
}

//from conn to hub
func (c *Client) readPump() {
	defer func() {
		c.hub.UnRegister() <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("error: ", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.Broadcast() <- message
	}
}

//from hub to conn
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				//hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
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

func SetHandlerWebSocket(rt *mux.Router) {
	rt.HandleFunc("/msg", wsMsgHandler)
	rt.HandleFunc("/mem", wsMemHandler)
}

/*
func serveWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func HandleWebSocket(rt *mux.Router) {
	hub := newHub()
	go hub.run()
	rt.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWebSocket(hub, w, r)
	})
}
*/
