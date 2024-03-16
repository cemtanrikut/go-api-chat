package app

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	conn     net.Conn
	username string
}

var clients []*Client

// ClientHandler defines the behavior of handling a client conn
type ClientHandler interface {
}

// ChatServer implements the ClientHandler interface
type ChatServer struct{}

func (cs *ChatServer) HandleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("New client connected:", conn.RemoteAddr())

	client := &Client{
		conn: conn,
	}

	clients = append(clients, client)

	// Asking username
	fmt.Fprint(conn, "Enter username: ")
	username, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	client.username = strings.TrimSpace(username)
	brodcastMessage(fmt.Sprintf("%s has joined the chat", client.username))

	//Handle incoming messages
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		if message == "/quit" {
			break
		}
		brodcastMessage(fmt.Sprintf("%s: %s", client.username, message))
	}

	// Remove client form the clients list
	for i, c := range clients {
		if c == client {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}

	brodcastMessage(fmt.Sprintf("%s has lef the chat", &client.username))

}

func brodcastMessage(message string) {
	fmt.Println("Brodcasting message:", message)
	for _, client := range clients {
		fmt.Fprintln(client.conn, message)
	}
}
