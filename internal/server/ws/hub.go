package ws

type NewSessionFunc func(*Client)

var hub *Hub

func StartHub(newSession NewSessionFunc) *Hub {
	hub = NewHub()
	go hub.run(newSession)

	return hub
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Broadcast() chan<- []byte {
	return h.broadcast
}

func (h *Hub) run(newSession NewSessionFunc) {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true

			c.OpenConnection()
			newSession(c)

		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				c.CancelSession()
				closeClientConnection(c, h)
			}

		case message := <-h.broadcast:
			for c := range h.clients {
				select {
				case c.ToClient() <- message:
				default:
					closeClientConnection(c, h)
				}
			}
		}
	}
}

func closeClientConnection(c *Client, h *Hub) {
	c.CloseConnection()
	delete(h.clients, c)
}
