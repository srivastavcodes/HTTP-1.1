package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"svr/internal/headers"
	"svr/internal/request"
	"svr/internal/response"
	"svr/internal/server"
	"syscall"
	"time"
)

const port = 7714

func hashString(bytes []byte) string {
	var str string
	for _, b := range bytes {
		str += fmt.Sprintf("%02x", b)
	}
	return str
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	svr, err := server.Serve(port, func(w *response.Writer, req *request.Request) {
		hdr := response.GetDefaultHeaders(0)

		code := response.StatusOk
		body := respond200()

		switch {
		case req.RequestLine.RequestTarget == "/yourproblem":
			body = respond400()
			code = response.StatusBadRequest

		case req.RequestLine.RequestTarget == "/myproblem":
			body = respond500()
			code = response.StatusInternalServerError

		case strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/"):
			reqTarget := req.RequestLine.RequestTarget

			res, err := http.Get("https://httpbin.org/" + reqTarget[len("/httpbin/"):])
			if err != nil {
				body = respond500()
				code = response.StatusInternalServerError
				break
			}
			_ = w.WriteStatusLine(response.StatusOk)

			hdr.Delete("Content-Length")
			hdr.Set("Transfer-Encoding", "chunked")
			hdr.Replace("Content-Type", "text/plain")
			hdr.Set("Trailer", "X-Content-SHA256")
			hdr.Set("Trailer", "X-Content-Length")

			_ = w.WriteHeaders(hdr)

			var compBody []byte
			for {
				data := make([]byte, 32)
				n, err := res.Body.Read(data)
				if err != nil {
					break
				}
				compBody = append(compBody, data[:n]...)

				_ = w.WriteBody([]byte(fmt.Sprintf("%x\r\n", n)))
				_ = w.WriteBody(data[:n])
				_ = w.WriteBody([]byte("\r\n"))
			}
			_ = w.WriteBody([]byte("0\r\n"))

			trailers := headers.NewHeaders()
			out := sha256.Sum256(compBody)
			trailers.Set("X-Content-SHA256", hashString(out[:]))
			trailers.Set("X-Content-Length", strconv.Itoa(len(compBody)))

			_ = w.WriteHeaders(trailers)
			return
		}
		hdr.Replace("Content-Length", strconv.Itoa(len(body)))
		hdr.Replace("Content-Type", "text/html")

		_ = w.WriteStatusLine(code)
		_ = w.WriteHeaders(hdr)
		_ = w.WriteBody(body)
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
