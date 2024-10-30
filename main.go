package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

// create custom types of messages
type MessageType int

const (
	WelcomeMessage MessageType = iota
	ClientListMessage
	ChatMessage
	SelectedClient
)

// Message represents a message with a type, payload, and timestamp
type Message struct {
	Type      MessageType `json:"type"`
	Sender    string      `json:"sender"`
	Relay     string      `json:"relay"`
	Payload   interface{} `json:"payload"` // use interface to store different kinds of data types
	Timestamp string      `json:"timestamp"`
}

// c2c relay
type HalfDuplex struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
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
	listener, err := net.Listen("tcp", ":9999")
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

// p2p: relay message {sender, receiver, payload}
func messageRelay(conn net.Conn, receiver string, msg string, now string) {
	// server relays the msg of net.conn
	relayPayload := Message{
		Type:      ChatMessage,
		Relay:     receiver,
		Sender:    conn.RemoteAddr().String(),
		Payload:   msg,
		Timestamp: now,
	}
	// serialize
	data, err := json.Marshal(relayPayload)
	if err != nil {
		log.Print("Error marshalling data: ", err)
		return
	}
	// get length
	length := int32(len(data))
	err = binary.Write(conn, binary.BigEndian, length)
	if err != nil {
		log.Print("Error writing data length: ", err)
		return
	}
	// write to receiver
	var thisReceiver net.Conn
	for _, client := range clients {
		if client.conn.RemoteAddr().String() == receiver {
			thisReceiver = client.conn
			break
		}
	}
	_, err = thisReceiver.Write(data)
	if err != nil {
		log.Print("Error writing data: ", err)
		return
	}
}

// sendClientList sends the list of available clients to each client
func sendClientList() {
	var clientAddrs []string
	for _, client := range clients {
		clientAddrs = append(clientAddrs, client.addr)
	}

	clientListMessage := Message{
		Type:      ClientListMessage,
		Payload:   clientAddrs,
		Timestamp: time.Now().Format("02-01-2006 15:04:05"),
	}

	broadcast(clientListMessage)
}

// handleClientConnection handles communication with a client
func handleClientConnection(conn net.Conn) {
	defer conn.Close()

	welcomeMessage := Message{
		Type:      WelcomeMessage,
		Sender:    conn.RemoteAddr().String(),
		Payload:   "WELCOME TO CLICHAT!",
		Timestamp: time.Now().Format("02-01-2006 15:04:05"),
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
			log.Printf("message from %s to %s: %s\n", conn.RemoteAddr(), msg.Relay, msg.Payload)
			msg.Timestamp = time.Now().Format("02-01-2006 15:04:05")
			//broadcast(msg)
			messageRelay(conn, msg.Relay, msg.Payload.(string), msg.Timestamp)
		case ClientListMessage:
			// Handle client list message if needed
		}
	}
}

