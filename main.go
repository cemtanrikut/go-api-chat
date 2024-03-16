package main

import (
	"fmt"
	"net"

	app "github.com/cemtanrikut/go-api-chat/app"
)

func main() {
	fmt.Println("Chat server started.")
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
