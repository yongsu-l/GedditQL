package main

import (
	"GedditQL/server"
	"GedditQL/server/interpreter"
	"GedditQL/server/parser"
	"log"
)

func main() {
	server := tcpserver.New("localhost:9999")

	server.OnNewClient(func(c *tcpserver.Client) {
		log.Println("New connection established")
		c.Send("Hello, welcome to GedditQL\n")
	})

	server.OnNewMessage(func(c *tcpserver.Client, query string) {
		// fmt.Print(message)
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		if r, err := parser.Tokenize(query); err != nil {
			c.Send("ERROR in tokenize\n")
		} else {
			// If there is an error with syntax, send the error to client
			if err := interpreter.CheckSyntax(r); err != nil {
				c.Send(err.Error() + "\n")
			} else {
				c.Send("Success\n")
				log.Println(interpreter.DescribeSelect(r))
			}
		}
	})

	server.OnClientConnectionClosed(func(c *tcpserver.Client, err error) {
		log.Println("Connection closed")
	})

	server.Listen()
}
