package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func read(conn net.Conn, alive chan bool) {
	for {

		reader := bufio.NewReader(conn)
		msg, err := reader.ReadString('\n')
		if err != nil {
			alive <- false
		}
		msg = strings.TrimSpace(msg)
		fmt.Printf(msg)

	}
}

func write(conn net.Conn) {
	stdin := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Enter text->")
		msg, _ := stdin.ReadString('\n')
		fmt.Fprintln(conn, msg)

	}
}
func main() {
	alive := make(chan bool)
	conn, _ := net.Dial("tcp", "127.0.0.1:8030")
	go write(conn)
	go read(conn, alive)
	//TODO Try to connect to the server
	//TODO Start asynchronously reading and displaying messages
	//TODO Start getting and sending user messages.
	for {
		if !<-alive {
			return
		}
	}

}
