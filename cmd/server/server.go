package main

import (
	"io"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"svr/internal/request"
	"svr/internal/response"
	"svr/internal/server"
	"syscall"
	"time"
)

const port = 7714

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	svr, err := server.Serve(port, func(w io.Writer, req *request.Request) *server.HandlerError {
		switch {
		case req.RequestLine.RequestTarget == "/yourproblem":
			return &server.HandlerError{
				StatusCode: response.StatusBadRequest,
				Message:    "your problem is not my problem\n",
			}
		case req.RequestLine.RequestTarget == "/myproblem":
			return &server.HandlerError{
				StatusCode: response.StatusInternalServerError,
				Message:    "Woopsie, my bad\n",
			}
		default:
			w.Write([]byte("All good, frfr\n"))
			return nil
		}
	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer svr.Close()
	logger.Info("Server started on", "port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	slog.Info("Server stopped gracefully", "Time", time.Now().Format("03:04:04"))
}
