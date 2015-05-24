package readraptor

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-martini/martini"
	"github.com/gorilla/websocket"
	"github.com/technoweenie/grohl"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// wsconn is an middleman between the websocket connection and the hub.
type wsconn struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// Account Id of the connection
	aid int64
}

func (c *wsconn) processMessage(m map[string]interface{}) {
	fmt.Println("msg", m)
	switch {
	case m["subscribe"] != "":
		reader, err := FindReaderByAccountIdDistinctId(c.aid, m["subscribe"].(string))
		if err != nil {
			grohl.Log(grohl.Data{
				"ws":    "FindReader",
				"error": err.Error(),
			})
			return
		}
		hub.subscribe <- &Channel{Conn: c, Id: reader.Id}
		fmt.Println("subscribe", reader.Id)
	}
}

// readPump pumps messages from the websocket connection to the hub.
func (c *wsconn) readPump() {
	defer func() {
		hub.unregister <- c
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}

		var j map[string]interface{}
		if err := json.Unmarshal(message, &j); err != nil {
			grohl.Log(grohl.Data{
				"ws":    "json parse error",
				"error": err.Error(),
			})
		} else {
			c.processMessage(j)
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *wsconn) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// write writes a message with the given message type and payload.
func (c *wsconn) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func WebsocketHandler(params martini.Params, w http.ResponseWriter, r *http.Request) {
	account, err := FindAccountByPublicKey(params["public_key"])
	if err != nil {
		log.Println(err)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	c := &wsconn{send: make(chan []byte, 256), ws: ws, aid: account.Id}
	hub.register <- c
	go c.writePump()
	c.readPump()
}
