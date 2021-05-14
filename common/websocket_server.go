package common

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

func (h *Hub) run() {
	for {
		select {
		case newClient := <-h.register:
			h.clients[newClient] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				h.clearClient(client)
				if len(h.clients) == 0 {
					return
				}
			}
		case message := <-h.broadcast:
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
