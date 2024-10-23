package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

// client ids
type Client struct {
	conn net.Conn
	addr string // use pure addr for now
}

var clients []Client

func main() {
	server() // starting websocket server
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
		client := Client{
			conn: conn,
			addr: conn.RemoteAddr().String(),
		}
		clients = append(clients, client)
		for _, client := range clients {
			fmt.Printf("Client %s connected\n", client.addr)
		}

		// handle multi clients connection
		go handleClientConnection(conn)
	}
}

// remove client id
func removeClientAddr(clientAddr string) {
	for i, client := range clients {
		if client.addr == clientAddr {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

// handle clients chosen comm
func commSwitch(conn net.Conn) {
	var availClients []string
	for _, client := range clients {
		availClients = append(availClients, client.addr)
	}

	// serialize availclients
	clientsAddrJSON, err := json.Marshal(availClients)
	if err != nil {
		log.Fatal("Error marshalling data: ", err)
		return
	}

	// write json data to connected client
	_, err = conn.Write(clientsAddrJSON)
	if err != nil {
		log.Fatal("Error writing to client: ", err)
		return
	}
}

// comm broadcasting
func broadcast(msg string) {
	message := []byte(msg)
	for _, client := range clients {
		_, err := client.conn.Write(message)
		if err != nil {
			log.Print("Error writing data: ", err)
		}
	}
}

// handle client connections
func handleClientConnection(conn net.Conn) {
	defer conn.Close() // close client connection if err occurs

	// each client shall receive a welcoming message
	welcomeMessage := []byte("WELCOME TO CLICHAT!\n")
	_, err := conn.Write(welcomeMessage)
	if err != nil {
		log.Println("Error sending welcome message:", err)
		return
	}

	// send connected clients to each client for them to query the comm protocol
	commSwitch(conn)

	// handle client comm protocol
	// client message transfer
	for {
		buffer := make([]byte, 1024) // increase buffer size for larger msg
		n, err := conn.Read(buffer)
		if err != nil {
			// client disconnection
			if err.Error() == "EOF" {
				log.Printf("%s has disconnected\n", conn.RemoteAddr())
				// remove client from ids list
				removeClientAddr(conn.RemoteAddr().String())
				return
			}
			log.Print("CLICHAT Server Error: ", err)
			return
		}

		// client message
		log.Printf("message from %s: %s\n", conn.RemoteAddr(), string(buffer[:n]))
		// broadcast this message
		broadcast(string(buffer[:n]))
	}
}
