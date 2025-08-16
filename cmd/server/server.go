package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"svr/internal/server"
	"syscall"
	"time"
)

const port = 7714

func main() {
	svr, err := server.Serve(port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer svr.Close()
	slog.Info("Server started on", "port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	slog.Info("Server stopped gracefully", "Time", time.Now().Format("03:04:04"))
}
