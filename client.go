package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	log.Print("client starting...")
	var err error

	node := flag.String("node", "127.0.0.1:2222", "ip address")
	flag.Parse()
	conn, err := net.Dial("tcp", *node)
	if err != nil {
		log.Fatal(err)
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Command: ")
		command, _ := reader.ReadString('\n')
		_, err = fmt.Fprintf(conn, command+"\n")
		if err != nil {
			log.Fatal(err)
		}
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print("Response: " + message)
	}
}
