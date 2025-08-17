package server

import (
	"fmt"
	"io"
	"net"
	"svr/internal/response"
)

type Server struct {
	closed bool
}

func runConnection(_ *Server, conn io.ReadWriteCloser) {
	defer conn.Close()
	headers := response.GetDefaultHeaders(0)
	_ = response.WriteStatusLine(conn, response.StatusOk)
	_ = response.WriteHeaders(conn, headers)
}

func (svr *Server) Close() error {
	svr.closed = true
	return nil
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
