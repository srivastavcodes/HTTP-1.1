package main

import (
	"log"
	"os"
	"os/signal"
	"sl/internal/server"
	"syscall"
)

func main() {
	srv := server.NewTcpServer("localhost", "4000")
	go func() {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigch
		log.Printf("recieved signal %s, shutting down...", sig)
		os.Exit(0)
	}()
	log.Fatal(srv.ServeTCP())
}
