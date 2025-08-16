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
		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("ERROR: ", err)
		}
		fmt.Printf("Request Line:\n")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)

		fmt.Printf("Headers:\n")
		req.Headers.ForEach(func(n, v string) {
			fmt.Printf("- %s: %s\n", n, v)
		})

		fmt.Printf("Body:\n")
		fmt.Printf("%s\n", req.Body)
	}
}
