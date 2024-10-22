package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func handle_incoming_messsage(conn net.Conn) {
	defer conn.Close()

	for {
		byte := make([]byte, 24)
		_, err := conn.Read(byte)
		if err != nil {
			log.Fatal("Error reading incoming data: ", err)
		}

		log.Println("Server: ", string(byte))

	}
}

func main() {

	conn, err := net.Dial("tcp", "localhost:80")
	if err != nil {
		log.Fatal("Error connecting to CLICHAT server: ", err)
	}

	byte_conn := []byte("client 1 connected")
	_, err = conn.Write(byte_conn)
	if err != nil {
		log.Fatal(err)
	}

	go handle_incoming_messsage(conn)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		fmt.Print("Enter message: ")
		message := scanner.Text()
		_, err := conn.Write([]byte(message))
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}
	}
}
