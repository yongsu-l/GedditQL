// Package tcpserver provides a tcpserver to accept queries
package tcpserver

import (
	"GedditQL/server/storage"
	"bufio"
	"encoding/gob"
	"log"
	"net"
)

// Client of each connections
type Client struct {
	conn   net.Conn
	Server *Server
}

// Server information
type Server struct {
	address                  string
	onNewClientCallback      func(c *Client)
	onClientConnectionClosed func(c *Client, err error)
	onNewMessage             func(c *Client, message string)
}

// Listen indefinitely until client closes connection
func (c *Client) Listen() {

	reader := bufio.NewReader(c.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			c.conn.Close()
			c.Server.onClientConnectionClosed(c, err)
			return
		}
		// log.Print(len(message))
		if len(message) > 1 {
			c.Server.onNewMessage(c, message)
		}
	}
}

// Send message back to client
func (c *Client) Send(res storage.Response) error {
	enc := gob.NewEncoder(c.Conn())
	return enc.Encode(&res)
}

// Conn establishes the connection
func (c *Client) Conn() net.Conn {
	return c.conn
}

// Close connection
func (c *Client) Close() error {
	return c.conn.Close()
}

// OnNewClient Called right after server starts listening new client
func (s *Server) OnNewClient(callback func(c *Client)) {
	s.onNewClientCallback = callback
}

// OnClientConnectionClosed Called right after connection closed
func (s *Server) OnClientConnectionClosed(callback func(c *Client, err error)) {
	s.onClientConnectionClosed = callback
}

// OnNewMessage Called when Client receives new message this function will handle all of the query
func (s *Server) OnNewMessage(callback func(c *Client, query string)) {
	s.onNewMessage = callback
}

// Listen starts the server listen to listen for incoming ocnnections
func (s *Server) Listen() {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal("Error starting TCP server.")
	}

	// Close connection once all other functions end
	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal("Error listening on connection")
		}
		client := &Client{
			conn:   conn,
			Server: s,
		}
		go client.Listen()
		s.onNewClientCallback(client)
	}
}

// New creates new tcp server instance
func New(address string) *Server {
	log.Println("Creating server with address", address)
	server := &Server{
		address: address,
	}

	server.OnNewClient(func(c *Client) {})
	server.OnNewMessage(func(c *Client, message string) {})
	server.OnClientConnectionClosed(func(c *Client, err error) {})

	return server
}
