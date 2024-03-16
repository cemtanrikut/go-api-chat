package app

import (
	"bufio"
	"net"
	"strings"
	"testing"
)

func TestHandleConnection(t *testing.T) {
	// Create an in-memory pipe
	serverConn, clientConn := net.Pipe()

	// Create a SimpleChatServer instance
	server := &ChatServer{}

	// Create a channel to receive the message from the server
	messageReceived := make(chan string)

	// Start a goroutine to read the message from the client connection
	go func() {
		defer clientConn.Close()
		scanner := bufio.NewScanner(clientConn)
		for scanner.Scan() {
			messageReceived <- scanner.Text()
		}
	}()

	// Start handling the connection
	go server.HandleConnection(serverConn)

	// Send a message from the client
	clientMessage := "Hello, Chat Server!"
	clientConn.Write([]byte(clientMessage + "\n"))

	// Wait for the message to be received by the server
	receivedMessage := <-messageReceived

	// Ensure that the server received the message correctly
	if !strings.Contains(receivedMessage, clientMessage) {
		t.Errorf("Expected server to receive message '%s', but received '%s'", clientMessage, receivedMessage)
	}
}
