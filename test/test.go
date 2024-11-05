package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:80")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	addr := conn.RemoteAddr().(*net.TCPAddr)
	hostnames, err := net.LookupAddr(addr.IP.String())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Hostname: %s\n", hostnames[0])
}
