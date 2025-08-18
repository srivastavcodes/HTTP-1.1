package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
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

	svr, err := server.Serve(port, func(w *response.Writer, req *request.Request) {
		hdr := response.GetDefaultHeaders(0)
		var err error

		code := response.StatusOk
		body := respond200()

		switch {
		case req.RequestLine.RequestTarget == "/yourproblem":
			body = respond400()
			code = response.StatusBadRequest

		case req.RequestLine.RequestTarget == "/myproblem":
			body = respond500()
			code = response.StatusInternalServerError
		}
		hdr.Replace("Content-Length", strconv.Itoa(len(body)))
		hdr.Replace("Content-Type", "text/html")

		err = w.WriteStatusLine(code)
		if err != nil {
			panic(fmt.Sprintf("shit blew up: %s", err))
		}
		err = w.WriteHeaders(hdr)
		if err != nil {
			panic(fmt.Sprintf("shit blew up: %s", err))
		}
		err = w.WriteBody(body)
		if err != nil {
			panic(fmt.Sprintf("shit blew up: %s", err))
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
