package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

// connection to clichat
var conn net.Conn

// handler: server incoming connection
func handleIncomingMessage(conn net.Conn) {
	defer conn.Close()

	// read the server message
	for {
		byte := make([]byte, 2048)
		n, err := conn.Read(byte)
		if err != nil {
			log.Println("Error reading incoming data: ", err)
			return
		}
		log.Printf("%s\n", string(byte[:n]))
	}
}

// connect to clichat server
func connectToServer() {
	log.Print("Starting Connection to CLICHAT")
	var err error
	conn, err = net.Dial("tcp", "localhost:80")
	if err != nil {
		log.Fatalf("Error connecting to CLICHAT server: %v", err)
	}

	byteConn := []byte("client 1 connected")
	_, err = conn.Write(byteConn)
	if err != nil {
		log.Fatalf("Error sending initial message: %v", err)
	}

	go handleIncomingMessage(conn)
}

// handle message relay
// relaying will use stdin for now
func messageRelay() {
	log.Print("Start Messaging")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Print("Message: ")
		message := scanner.Text()
		_, err := conn.Write([]byte(message))
		if err != nil {
			log.Println("Error sending message: ", err)
			return
		}
		fmt.Print("\n")
	}
}

func main() {
	connectToServer() // start clichat connection
	messageRelay()    // start relaying message
}
