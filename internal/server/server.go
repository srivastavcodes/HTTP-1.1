package server

import (
	"errors"
	"fmt"
	"io"
	"net"
	"svr/internal/request"
	"svr/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	closed  bool
	handler Handler
}

func runConnection(svr *Server, conn io.ReadWriteCloser) {
	defer conn.Close()
	rw := response.NewWriter(conn)

	req, err := request.RequestFromReader(conn)
	if err != nil {
		_ = rw.WriteStatusLine(response.StatusBadRequest)
		return
	}
	svr.handler(rw, req)
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

func Serve(port uint16, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, errors.New("error starting server")
	}
	server := &Server{
		closed:  false,
		handler: handler,
	}
	go runServer(server, listener)
	return server, nil
}
