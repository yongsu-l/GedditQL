package main

import (
	"GedditQL/client/parse"
	"GedditQL/server/storage"
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":9999")
	defer conn.Close()

	if err != nil {
		log.Fatalln(err)
	}

	// decoder.Decode(res)
	// log.Println(res)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		query, err := reader.ReadString('\n')
		decoder := gob.NewDecoder(conn)

		if err != nil {
			log.Fatalln(err.Error())
		} else if query == "quit\n" {
			fmt.Println("Exiting...")
			os.Exit(0)
		}

		// Send to query
		_, err = fmt.Fprintf(conn, query+"\n")

		if err != nil {
			log.Fatal(err.Error())
		}

		var res storage.Response
		// log.Print("Hang")
		err = decoder.Decode(&res)
		if err != nil {
			log.Fatal(err)
		}

		parse.Table(res)
	}

}
