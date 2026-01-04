package ws

type Message struct {
	ChatroomID uint
	Data       []byte
}

type Hub struct {
	Rooms      map[uint]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[uint]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.Register:
			if h.Rooms[client.ChatroomID] == nil {
				h.Rooms[client.ChatroomID] = make(map[*Client]bool)
			}
			h.Rooms[client.ChatroomID][client] = true

		case client := <-h.Unregister:
			if clients, ok := h.Rooms[client.ChatroomID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
				}
				if len(clients) == 0 {
					delete(h.Rooms, client.ChatroomID)
				}
			}

		case msg := <-h.Broadcast:
			if clients, ok := h.Rooms[msg.ChatroomID]; ok {
				for client := range clients {
					select {
					case client.Send <- msg.Data:
					default:
						close(client.Send)
						delete(clients, client)
					}
				}
			}
		}
	}
}
