package websocket

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
			ncMsg := jsonMustMarshal(newRegisterMessage(newClient.name))
			for client := range h.clients {
				select {
				case client.send <- ncMsg:
				default:
					h.clearClient(client)
				}
				cMsg := jsonMustMarshal(newRegisterMessage(client.name))
				select {
				case newClient.send <- cMsg:
				default:
					h.clearClient(client)
				}
			}
			newClient.send <- ncMsg
			h.clients[newClient] = true
		case tClient := <-h.unregister:
			if _, ok := h.clients[tClient]; ok {
				tcMsg := jsonMustMarshal(newUnregisterMessage(tClient.name))
				for c := range h.clients {
					select {
					case c.send <- tcMsg:
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
	close(client.send)
	delete(h.clients, client)
}

func (h *Hub) Register() chan<- *Client {
	return h.register
}

func (h *Hub) Unregister() chan<- *Client {
	return h.unregister
}

func (h *Hub) Broadcast() chan<- []byte {
	return h.broadcast
}
