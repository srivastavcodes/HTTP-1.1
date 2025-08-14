package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)

		var str string
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				break
			}
			data = data[:n]
			// prints the line only when '\n' is encountered
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				str += string(data[:i])
				data = data[i+1:]
				out <- str
				str = ""
			}
			str += string(data)
		}
		if len(str) != 0 {
			out <- str
		}
	}()
	return out
}

func main() {
	listener, err := net.Listen("tcp", ":7714")
	fmt.Printf("TCP listener started\n")
	if err != nil {
		log.Fatal("ERROR: ", err.Error())
	}
	for {
		listen, err := listener.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		}
		for line := range getLinesChannel(listen) {
			fmt.Printf("read: %s\n", line)
		}
	}
}
