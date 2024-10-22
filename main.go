package main

import (
	"fmt"
	"log"
	"net"
)

// client ids
var client_ids []string

// remove client id
func remove_client_id(client_id string) {
	for i, id := range client_ids {
		if id == client_id {
			client_ids = append(client_ids[:i], client_ids[i+1:]...)
			break
		}
	}
}

// CLICHAT Server
func server() {

	// set up listener for port :80
	// localhost: simulated client-to-client will have to bind to :80
	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatal(err)
	}
	log.Print("CLICHAT Starting")

	// handle incoming client connections
	// store each connection as ids in a slice
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print("Error accepting client connections: ", err)
		}

		// conn established with conn.RemoteAddr
		client_ids = append(client_ids, conn.RemoteAddr().String())
		for _, addr := range client_ids {
			fmt.Printf("Client %s connected\n", addr)
		}

		// handle multi client connection
		go handleClientConnection(conn)
	}

}

func main() {

	server() // staring websocket server

}

// handle client connections
func handleClientConnection(conn net.Conn) {
	defer conn.Close() // close client connection if err occurs

	// each client shall recieve a welcoming message
	welcomeMessage := []byte("WELCOME TO CLICHAT!")
	_, err := conn.Write(welcomeMessage)
	if err != nil {
		log.Println("Error sending welcome message:", err)
		return
	}

	// handle client comm protocol
	// client message transfer
	for {
		byte := make([]byte, 24) // increase buffer size for larger msg
		_, err := conn.Read(byte)
		if err != nil {
			// client disconnection
			if err.Error() == "EOF" {
				log.Printf("%s has diconnected", conn.RemoteAddr())
				// remove client from ids list
				remove_client_id(conn.RemoteAddr().String())
				for _, addr := range client_ids {
					fmt.Printf("Client %s connected\n", addr)
				}
				return
			}
			log.Print("CLICHAT Server Error: ", err)
			return
		}

		// client message
		log.Printf("message from %s : %s", conn.RemoteAddr(), string(byte))

	}
}
