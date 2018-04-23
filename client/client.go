package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {

	// Connect to the socket
	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println("Error establishing connection:", err.Error())
		os.Exit(1)
	}
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		os.Exit(1)
	}
	fmt.Println(message)

	for {
		// Read input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		query, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading string:", err.Error())
			os.Exit(1)
		}
		if query == "quit\n" {
			fmt.Println("Quitting...")
			os.Exit(0)
		}
		// Send to db server
		fmt.Fprintf(conn, query+"\n")
		// Listen for reply
		message, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil {
			fmt.Println("Error reading buffer:", err.Error())
			os.Exit(1)
		}

		fmt.Println(message)
	}
}
