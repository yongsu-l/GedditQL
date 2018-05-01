package main

import (
	"GedditQL/server"
	// "GedditQL/server/interpreter"
	"GedditQL/server/linter"
	"GedditQL/server/parser"
	"GedditQL/server/storage"
	// "encoding/gob"
	"log"
	"strings"
	// "reflect"
)

func main() {
	// Start logging based on line
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	server := tcpserver.New("localhost:9999")

	Linter := linter.NewLinter("./server/grammar.txt")

	// Initialize a new db
	db, _ := storage.New("db", "store")

	server.OnNewClient(func(c *tcpserver.Client) {
		log.Println("New connection established")

		// res := storage.Response{Msg: "Hello, welcome to GedditQL"}

		// c.Send(res)
	})

	server.OnNewMessage(func(c *tcpserver.Client, query string) {
		log.Print(query)

		if chk := Linter(strings.TrimSpace(query), "query"); chk {
			// If query has valid syntax, tokenize the evaluate
			log.Print("WORKS")
			if r, err := parser.Tokenize(query); err != nil {
				res := storage.Response{Err: err.Error()}
				c.Send(res)
			} else {
				res, err := db.EvaluateQuery(r)
				if err = c.Send(res); err != nil {
					log.Fatal(err)
				} else {
					// No error
				}
			}
		} else {
			log.Print(chk)
			res := storage.Response{Err: "Invalid Syntax"}
			c.Send(res)
		}
	})

	server.OnClientConnectionClosed(func(c *tcpserver.Client, err error) {
		log.Println("Connection closed")
	})

	server.Listen()
}
