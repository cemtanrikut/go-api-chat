package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/cemtanrikut/go-api-chat/app"
)

// ClientHandler defines the behavior of handling a client connection.
type ClientHandler interface {
	HandleConnection(conn net.Conn)
}

// SimpleChatServer implements the ClientHandler interface.
type SimpleChatServer struct{}

func (s *SimpleChatServer) HandleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("New client connected:", conn.RemoteAddr())

	client := &Client{
		conn: conn,
	}

	clients = append(clients, client)

	// Ask for username
	fmt.Fprint(conn, "Enter your username: ")
	username, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	client.username = strings.TrimSpace(username)
	broadcastMessage(fmt.Sprintf("%s has joined the chat", client.username))

	// Handle incoming messages
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		if message == "/quit" {
			break
		}
		broadcastMessage(fmt.Sprintf("%s: %s", client.username, message))
	}

	// Remove client from the list
	for i, c := range clients {
		if c == client {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}

	broadcastMessage(fmt.Sprintf("%s has left the chat", client.username))
}

// Client represents a connected client.
type Client struct {
	conn     net.Conn
	username string
}

var clients []*Client

func main() {
	fmt.Println("Chat Server Started")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()

	handler := &app.ChatServer{}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handler.HandleConnection(conn)
	}
}

func broadcastMessage(message string) {
	fmt.Println("Broadcasting message:", message)
	for _, client := range clients {
		fmt.Fprintln(client.conn, message)
	}
}
