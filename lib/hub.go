package readraptor

import "fmt"

// Hub maintains the set of active websocket connections
type Hub struct {
	// Registered connections.
	connections map[*wsconn]bool

	subscriptions map[int64]*wsconn

	// Inbound messages from the connections.
	broadcast chan *Broadcast

	// Register requests from the connections.
	register chan *wsconn

	// Unregister requests from connections.
	unregister chan *wsconn

	// Register distinct id listening channel
	subscribe chan *Channel

	// Unregister distinct id listening channel
	unsubscribe chan *Channel
}

type Broadcast struct {
	ConnIds []int64
	Message []byte
}

type Channel struct {
	Conn *wsconn
	Id   int64
}

var hub = Hub{
	broadcast:   make(chan *Broadcast),
	register:    make(chan *wsconn),
	subscribe:   make(chan *Channel),
	unsubscribe: make(chan *Channel),
	unregister:  make(chan *wsconn),

	connections:   make(map[*wsconn]bool),
	subscriptions: make(map[int64]*wsconn),
}

func (h *Hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)

				for k, v := range h.subscriptions {
					if v == c {
						fmt.Println("removing", k)
						delete(h.subscriptions, k)
					}
				}
			}
		case c := <-h.subscribe:
			h.subscriptions[c.Id] = c.Conn
		case c := <-h.unsubscribe:
			if _, ok := h.subscriptions[c.Id]; ok {
				fmt.Println("unsubscribe", c.Id)
				delete(h.subscriptions, c.Id)
			}
		case b := <-h.broadcast:
			for _, id := range b.ConnIds {
				if c, ok := h.subscriptions[id]; ok {
					select {
					case c.send <- b.Message:
					default:
						close(c.send)
						delete(h.connections, c)
					}
				}
			}
		}
	}
}
