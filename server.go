// Server to host GedditQL Database. Will default to port 3306

package main

import (
	"fmt"
	"net"
	"os"
)

const (
	// ConnHost is for hosting localhost
	ConnHost = "localhost"
	// ConnPort is the port for the database server
	ConnPort = "3306"
	// ConnType is the type of connection that the db server will accept
	ConnType = "tcp"
)

func main() {
	// Listen for incoming connections
	l, err := net.Listen(ConnType, ConnHost+":"+ConnPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close listener when application closes
	defer l.Close()

	fmt.Println("Listening on " + ConnHost + ":" + ConnPort)
	for {
		// Listen for incoming connections
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	// make a buffer to hold all of the incoming data
	buf := make([]byte, 1024)
	// Read the incoming data to buffer
	reqLen, err := conn.Read(buf)
	fmt.Println(reqLen)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		os.Exit(1)
	}
	conn.Write([]byte("Message received."))
	conn.Close()
}
