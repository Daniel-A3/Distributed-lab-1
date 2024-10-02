package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strings"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	// TODO: all
	// Deal with an error event.
	fmt.Printf(err.Error())
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			handleError(err)
			continue
		}
		conns <- conn
	}
}

func handleClient(client net.Conn, clientid int, msgs chan Message, clients map[int]net.Conn) {
	// TODO: all
	// So long as this connection is alive:
	// Read in new messages as delimited by '\n's
	// Tidy up each message and add it to the messages channel,
	// recording which client it came from.

	for {

		reader := bufio.NewReader(client)
		message, err := reader.ReadString('\n')
		message = strings.TrimSpace(message) //may not need

		if err != nil {
			handleError(err)
			delete(clients, clientid)
			fmt.Printf("connection closed for client: " + string(clientid))
			return
		}

		msg := Message{
			sender:  clientid,
			message: message,
		}

		msgs <- msg
	}

}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	//TODO Create a Listener for TCP connections on the port given above.

	ln, err := net.Listen("tcp", *portPtr)
	if err != nil {
		handleError(err)
		return
	}
	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)

	//Start accepting connections
	go acceptConns(ln, conns)
	clientID := 0
	for {
		select {
		case conn := <-conns:
			//TODO Deal with a new connection
			// - assign a client ID
			clients[clientID] = conn
			go handleClient(conn, clientID, msgs, clients)
			clientID++
			// - add the client to the clients map
			// - start to asynchronously handle messages from this client
		case msg := <-msgs:
			//TODO Deal with a new message
			// Send the message to all clients that aren't the sender
			for ID, client := range clients {
				if msg.sender != ID {
					fmt.Fprintln(client, msg.message)
				}
			}

		}
	}
}
