package server

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"svr/internal/request"
	"svr/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

type Server struct {
	closed  bool
	handler Handler
}

func runConnection(svr *Server, conn io.ReadWriteCloser) {
	defer conn.Close()
	headers := response.GetDefaultHeaders(0)

	req, err := request.RequestFromReader(conn)
	if err != nil {
		_ = response.WriteStatusLine(conn, response.StatusBadRequest)
		_ = response.WriteHeaders(conn, headers)
		return
	}
	writer := bytes.NewBuffer(make([]byte, 0))

	var body []byte
	var code = response.StatusOk

	handlerErr := svr.handler(writer, req)
	if handlerErr != nil {
		body = []byte(handlerErr.Message)
		code = handlerErr.StatusCode
	} else {
		body = writer.Bytes()
	}
	headers.Replace("Content-Length", strconv.Itoa(len(body)))

	_ = response.WriteStatusLine(conn, code)
	_ = response.WriteHeaders(conn, headers)

	conn.Write(body)
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
		return nil, fmt.Errorf("error starting server")
	}
	server := &Server{
		closed:  false,
		handler: handler,
	}
	go runServer(server, listener)
	return server, nil
}
