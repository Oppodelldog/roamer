package ws

import (
	"log"
	"net/http"
)

const channelBuffer = 256

// ServeWs handles websocket requests from the peer
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)

		return
	}

	sendCh := make(chan []byte, channelBuffer)
	respCh := make(chan []byte, channelBuffer)
	client := &Client{hub: hub, conn: conn, toClient: sendCh, toServer: respCh}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
