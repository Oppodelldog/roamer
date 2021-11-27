package ws

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gorilla/websocket"

	"io"
	"net/http"
	"sync"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 20 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024 * 1024
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

const readBuffer = 1024
const writeBuffer = 1024

var upgrader = websocket.Upgrader{
	ReadBufferSize:  readBuffer,
	WriteBufferSize: writeBuffer,
	CheckOrigin: func(r *http.Request) bool {
		for k, v := range r.Header {
			fmt.Println(k, v)
		}

		return true
	},
}

func (c *Client) Send(b []byte) {
	if !c.IsActive() {
		return
	}
	c.toClient <- b
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	CancelSession context.CancelFunc

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	toClient     chan []byte
	toServer     chan []byte
	toServerLock sync.Mutex

	active     bool
	activeLock sync.Mutex
}

func (c *Client) ToClient() chan<- []byte {
	return c.toClient
}

func (c *Client) ToServer() <-chan []byte {
	return c.toServer
}

func (c *Client) GetAddress() string {
	return c.conn.RemoteAddr().String()
}

func (c *Client) OpenConnection() {
	c.setActive(true)
}

func (c *Client) CloseConnection() {
	c.setActive(false)
	close(c.toClient)
}

func (c *Client) IsActive() bool {
	c.activeLock.Lock()
	defer c.activeLock.Unlock()

	return c.active
}

func (c *Client) setActive(value bool) {
	c.activeLock.Lock()
	defer c.activeLock.Unlock()

	c.active = value
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c

		err := c.conn.Close()
		if err != nil {
			fmt.Printf("error closing connection: %v\n", err)
		}
	}()

	c.conn.SetReadLimit(maxMessageSize)

	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		panic(err)
	}

	c.conn.SetPongHandler(func(string) error { return c.conn.SetReadDeadline(time.Now().Add(pongWait)) })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
				websocket.CloseNoStatusReceived,
			) {
				fmt.Printf("unexpected close error: %v", err)
				close(c.toServer)
			}

			break
		}

		message = bytes.TrimSpace(bytes.ReplaceAll(message, newline, space))

		c.toServerLock.Lock()

		if c.toServer != nil {
			c.toServer <- message
		} else {
			fmt.Printf("toServer Channel is nil, message will be ignored")
		}

		c.toServerLock.Unlock()
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer c.cleanupWritePump(ticker)

	for {
		select {
		case message, ok := <-c.toClient:
			setDeadline(c)

			if !ok {
				closeWebsocket(c)

				return
			}

			if err := write(c, message); err != nil {
				return
			}

		case <-ticker.C:
			if err := ping(c); err != nil {
				return
			}
		}
	}
}

func (c *Client) cleanupWritePump(ticker *time.Ticker) {
	ticker.Stop()

	err := c.conn.Close()
	if err != nil {
		fmt.Printf("error closing client connection: %v", err)
	}
}

func write(c *Client, message []byte) error {
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		fmt.Printf("error NextWriter %v", err)

		return err
	}

	_, err = w.Write(withNewLine(message))
	if err != nil {
		fmt.Printf("error writing message: %v", err)

		return err
	}

	n := len(c.toClient)
	for i := 0; i < n; i++ {
		_, err := w.Write(withNewLine(<-c.toClient))
		if err != nil {
			fmt.Printf("error writing c.toClient: %v", err)

			return err
		}
	}

	if err := w.Close(); err != nil && err != io.EOF {
		fmt.Printf("error closing writer: %v", err)

		return err
	}

	return nil
}

func setDeadline(c *Client) {
	err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err != nil {
		fmt.Printf("error setting write deadline: %v", err)
	}
}

func closeWebsocket(c *Client) {
	err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
	if err != nil {
		fmt.Printf("error writing CloseMessage: %v", err)
	}
}

func ping(c *Client) error {
	if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		fmt.Printf("Error set write deadline: %v", err)

		return err
	}

	if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
		fmt.Printf("Error writing ping message: %v", err)

		return err
	}

	return nil
}

func withNewLine(message []byte) []byte {
	return append(message, newline...)
}
