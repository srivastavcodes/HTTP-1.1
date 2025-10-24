package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sl/server/html"
	"strings"
	"time"
)

type HTTPServer struct {
	Host string
	Port string
	log  *log.Logger
}

func NewHttpServer(host string, port string) *HTTPServer {
	logger := log.New(os.Stdout, "HTTP :: ", log.LstdFlags|log.Lmsgprefix)
	return &HTTPServer{
		Host: host,
		Port: port,
		log:  logger,
	}
}

// ServeHTTP starts the HTTP server and listens for incoming connections. It
// accepts connections in a loop and spawns a goroutine to handle each one.
//
// Returns an error if the server fails to start listening on the configured
// address.
func (s *HTTPServer) ServeHTTP() error {
	address := fmt.Sprintf("%s:%s", s.Host, s.Port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen on port=%s. err=%s", s.Port, err)
	}
	defer listener.Close()
	s.logf("Server listening on port=%s", s.Port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logf("error accepting connection: %s\n", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

// handleConnection processes a single HTTP client connection. It reads the
// HTTP request line and headers, routes the request based on the path, and
// sends back an appropriate HTTP response before closing the connection.
func (s *HTTPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	s.logf("new connection from: %s\n", conn.RemoteAddr().String())
	conn.SetReadDeadline(time.Now().Add(15 * time.Second))

	// reader contains the http request
	reader := bufio.NewReader(conn)

	reqLine, err := reader.ReadString('\n')
	if err != nil {
		s.logf("error reading request line: %s\n", err)
		return
	}
	// request line = METHOD PATH VERSION\r\n
	rlSlice := strings.Fields(strings.TrimSpace(reqLine))
	if len(rlSlice) != 3 {
		s.logf("invalid request from %s: %s\n", conn.RemoteAddr(), reqLine)
		s.sendErrorResponse(conn, 400, "Bad Request")
		return
	}
	var (
		method  = rlSlice[0]
		path    = rlSlice[1]
		version = rlSlice[2]
	)
	s.logf("request received :: method=%s url=%s version=%s", method, path, version)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			s.logf("error reading headers: %s\n", err)
			return
		}
		if strings.TrimSpace(line) == "" {
			break
		}
		s.logf("headers received: %s", strings.TrimSpace(line))
	}
	response := []byte(reqLine)

	if strings.Contains(path, "hello") {
		response = SayHello().WriteResponse()
	} else {
		s.sendErrorResponse(conn, 400, "Bad Request")
	}
	_, err = conn.Write(response)
	if err != nil {
		s.logf("error sending response: %s\n", err)
		return
	}
	s.logf("response sent to client %s\n", conn.RemoteAddr().String())
}

// TODO: enhance the error handling and html with dynamic data

// sendErrorResponse sends an error response back to the client with a html
// body containing error details.
func (s *HTTPServer) sendErrorResponse(conn net.Conn, code int, text string) {
	res := NewHTTPResponse()

	res.SetStatus(code, text)

	res.SetHeader("Content-Type", "text/html; charset=utf-8")
	res.SetHeader("Connection", "close")

	res.SetBody(html.BadRequestHTML)

	conn.Write(res.WriteResponse())
}

// logf logs defensively
func (s *HTTPServer) logf(format string, args ...any) {
	if s.log != nil {
		s.log.Printf(format, args...)
	} else {
		logger := log.New(os.Stdout, "HTTP :: ", log.LstdFlags|log.Lmsgprefix)
		logger.Printf(format, args...)
	}
}
