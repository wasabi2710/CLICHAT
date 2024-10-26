package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

// create custom types of messages
type MessageType int

const (
	WelcomeMessage MessageType = iota
	ClientListMessage
	ChatMessage
)

// Message represents a message with a type and payload
type Message struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"` // use interface to store differen kind of data types
}

// Client represents a connected client
type Client struct {
	conn net.Conn
	addr string
}

var clients []Client

func main() {
	server()
}

// server starts the CLICHAT server
func server() {
	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatal(err)
	}
	log.Print("CLICHAT Starting")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print("Error accepting client connections: ", err)
			continue
		}

		client := Client{
			conn: conn,
			addr: conn.RemoteAddr().String(),
		}
		clients = append(clients, client)
		for _, client := range clients {
			fmt.Printf("Client %s connected\n", client.addr)
		}

		go handleClientConnection(conn)
	}
}

// removeClientAddr removes a client from the list by address
func removeClientAddr(clientAddr string) {
	for i, client := range clients {
		if client.addr == clientAddr {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

// sendMessage sends a message to a client
func sendMessage(conn net.Conn, msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	length := int32(len(data))
	err = binary.Write(conn, binary.BigEndian, length)
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}

// broadcast sends a message to all clients
func broadcast(msg Message) {
	for _, client := range clients {
		err := sendMessage(client.conn, msg)
		if err != nil {
			log.Print("Error writing data: ", err)
		}
	}
}

// sendClientList sends the list of available clients to each client
func sendClientList() {
	var clientAddrs []string
	for _, client := range clients {
		clientAddrs = append(clientAddrs, client.addr)
	}

	clientListMessage := Message{
		Type:    ClientListMessage,
		Payload: clientAddrs,
	}

	broadcast(clientListMessage)
}

// handleClientConnection handles communication with a client
func handleClientConnection(conn net.Conn) {
	defer conn.Close()

	welcomeMessage := Message{
		Type:    WelcomeMessage,
		Payload: "WELCOME TO CLICHAT!",
	}
	err := sendMessage(conn, welcomeMessage)
	if err != nil {
		log.Println("Error sending welcome message:", err)
		return
	}

	sendClientList()

	for {
		var length int32
		err := binary.Read(conn, binary.BigEndian, &length)
		if err != nil {
			log.Printf("%s has disconnected\n", conn.RemoteAddr())
			removeClientAddr(conn.RemoteAddr().String())
			sendClientList()
			return
		}

		buffer := make([]byte, length)
		_, err = conn.Read(buffer)
		if err != nil {
			log.Print("CLICHAT Server Error: ", err)
			return
		}

		var msg Message
		err = json.Unmarshal(buffer, &msg)
		if err != nil {
			log.Print("Error unmarshalling message: ", err)
			continue
		}

		switch msg.Type {
		case ChatMessage:
			log.Printf("message from %s: %s\n", conn.RemoteAddr(), msg.Payload)
			broadcast(msg)
		case ClientListMessage:
			// Handle client list message if needed
		}
	}
}
