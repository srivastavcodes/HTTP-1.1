package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

// TcpServer defines the parameters for running a tcp server. Must be
// initialized
type TcpServer struct {
	Host string
	Port string
	log  *log.Logger
}

func NewTcpServer(host string, port string) *TcpServer {
	logger := log.New(os.Stdout, "Tcp :: ", log.LstdFlags|log.Lmsgprefix)
	return &TcpServer{
		Host: host,
		Port: port,
		log:  logger,
	}
}

// ServeTCP starts a new TCP server listening on the configured Host and Port.
func (s *TcpServer) ServeTCP() error {
	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("error listening on address %s: %s", addr, err)
	}
	defer listener.Close()
	s.logf("listening on address: %s\n", addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logf("error accepting connection: %s\n", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *TcpServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	s.logf("new connection from: %s\n", conn.RemoteAddr().String())

	conn.SetDeadline(time.Now().Add(30 * time.Second))
	info := fmt.Sprintf("connected at %s\n", time.Now().Format(time.RFC822))

	_, err := conn.Write([]byte(info))
	if err != nil {
		s.logf("error sending response: %s\n", err)
		return
	}
read:
	for {
		buf := make([]byte, 1024)

		n, err := conn.Read(buf)
		if err != nil {
			s.logf("error reading response: %s\n", err)
			break
		}
		s.logf("read %d bytes: %s.", n, string(buf[:n-1]))

		// echo to the client
		_, err = conn.Write(buf[:n])
		if err != nil {
			s.logf("error sending response: %s\n", err)
			break read
		}
		conn.SetDeadline(time.Now().Add(30 * time.Second))
	}
	s.logf("connection with client=%s closed\n", conn.RemoteAddr().String())
}

func (s *TcpServer) logf(format string, args ...any) {
	if s.log != nil {
		s.log.Printf(format, args...)
	} else {
		logger := log.New(os.Stdout, "TcpServer", log.LstdFlags|log.Lshortfile)
		logger.Printf(format, args...)
	}
}
