package common

import "log"

type IHub interface {
	Run()
	Register() chan<- *Client
	UnRegister() chan<- *Client
	Broadcast() chan<- []byte
}
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case newClient := <-h.register:
			log.Println("msghub: register")
			h.clients[newClient] = true
		case client := <-h.unregister:
			log.Println("msghub: unregister")
			if _, ok := h.clients[client]; ok {
				h.clearClient(client)
				if len(h.clients) == 0 {
					return
				}
			}
		case message := <-h.broadcast:
			log.Println("msghub: broadcast")
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					h.clearClient(client)
				}
			}
		}
	}
}

func (h *Hub) clearClient(client *Client) {
	//close(client.sendMember)
	close(client.send)
	delete(h.clients, client)
}

func (h *Hub) Register() chan<- *Client {
	return h.register
}

func (h *Hub) UnRegister() chan<- *Client {
	return h.unregister
}

func (h *Hub) Broadcast() chan<- []byte {
	return h.broadcast
}

type MemHub struct {
	*Hub
}

func newMemHub() *MemHub {
	return &MemHub{newHub()}
}

func (h *MemHub) Run() {
	for {
		select {
		case newClient := <-h.register:
			log.Println("memhub: register")
			for client := range h.clients {
				select {
				case client.send <- []byte(newClient.name):
				default:
					h.clearClient(client)
				}
				select {
				case newClient.send <- []byte(client.name):
				default:
					h.clearClient(client)
				}
			}
			newClient.send <- []byte(newClient.name)
			h.clients[newClient] = true
		case tClient := <-h.unregister:
			log.Println("memhub: unregister")
			if _, ok := h.clients[tClient]; ok {
				for c := range h.clients {
					select {
					case c.send <- []byte(tClient.name):
					default:
						h.clearClient(c)
					}
				}
				h.clearClient(tClient)
				if len(h.clients) == 0 {
					return
				}
			}
		case message := <-h.broadcast:
			log.Println("memhub: broadcast")
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					h.clearClient(client)
				}
			}
		}
	}
}
