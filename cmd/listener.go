package main

import (
	"fmt"
	"log"
	"net"
	"svr/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":7714")
	fmt.Printf("TCP listener started\n")
	if err != nil {
		log.Fatal("ERROR: ", err.Error())
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("ERROR: ", err)
		}
		rl, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("ERROR: ", err)
		}
		fmt.Printf("Request Line:\n")
		fmt.Printf("- Method: %s\n", rl.RequestLine.Method)
		fmt.Printf("- Target: %s\n", rl.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", rl.RequestLine.HttpVersion)
	}
}
