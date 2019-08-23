package main

import (
	"log"

	"github.com/gorilla/websocket"
)

var clients []*Clients

type Client struct {
	conn *websocket.Conn
	send chan *Message

	roomID string
	user   *User
}

const messageBufferSize = 256

func newClient(conn *websocket.Conn, roomID string, u *User) {
	c := &Client{
		conn:   conn,
		send:   make(chan *Message, messageBufferSize),
		roomID: roomID,
		user:   u,
	}

	clients = append(clients, c)

	go c.readLoop()
	go c.writeLoop()
}

// Close client's connection and channel
func (c *Client) Close() {
	for i, client := range clients {
		if client == c {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}

	close(c.send)
	c.conn.Close()
	log.Printf("close connection. addr: %s", c.conn.RemoteAddr())
}

func (c *Client) readLoop() {
	for {
		m, err := c.read()
		if err != nil {
			log.Println("read message error: ", err)
			break
		}

		m.create()
	}
}
