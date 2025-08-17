package response

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"svr/internal/headers"
)

type Response struct {
}

type StatusCode int

const (
	StatusOk                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusNotFound            StatusCode = 404
	StatusInternalServerError StatusCode = 500
)

func GetDefaultHeaders(conLen int) *headers.Headers {
	hdr := headers.NewHeaders()

	hdr.Set("Content-Length", strconv.Itoa(conLen))
	hdr.Set("Connection", "close")
	hdr.Set("Content-Type", "text/plain")

	return hdr
}

type Writer struct {
	writer io.Writer
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{
		writer: writer,
	}
}

func (w *Writer) WriteStatusLine(code StatusCode) error {
	var line []byte

	switch code {
	case StatusOk:
		line = []byte("HTTP/1.1 200 Ok\r\n")

	case StatusBadRequest:
		line = []byte("HTTP/1.1 400 Bad Request\r\n")

	case StatusNotFound:
		line = []byte("HTTP/1.1 404 Not Found\r\n")

	case StatusInternalServerError:
		line = []byte("HTTP/1.1 500 Internal Server Error\r\n")
	default:
		return errors.New("unrecognized status code")
	}
	_, err := w.writer.Write(line)
	return err
}

func (w *Writer) WriteBody(b []byte) error {
	_, err := w.writer.Write(b)
	return err
}

func (w *Writer) WriteHeaders(hdr *headers.Headers) error {
	var pair []byte
	hdr.ForEach(func(n, v string) {
		pair = fmt.Appendf(pair, "%s: %s\r\n", n, v)
	})
	pair = fmt.Append(pair, "\r\n")

	_, err := w.writer.Write(pair)
	return err
}
