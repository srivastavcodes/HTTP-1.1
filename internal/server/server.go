package server

import (
	"fmt"
	"io"
	"net"
)

type Server struct {
	closed bool
}

func runConnection(svr *Server, conn io.ReadWriteCloser) {
	out := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello World!`")
	conn.Write(out)
	conn.Close()
}

func runServer(svr *Server, listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if svr.closed {
			return
		}
		if err != nil {
			return
		}
		go runConnection(svr, conn)
	}
}

func Serve(port uint16) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("error starting server")
	}
	server := &Server{
		closed: false,
	}
	go runServer(server, listener)
	return server, nil
}

func (svr *Server) Close() error {
	svr.closed = true
	return nil
}
