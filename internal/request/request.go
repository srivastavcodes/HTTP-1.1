package request

import (
	"bytes"
	"fmt"
	"io"
	"svr/internal/headers"
)

var (
	ERROR_BAD_START_LINE         = fmt.Errorf("malformed request-line")
	ERROR_REQUEST_IN_ERROR_STATE = fmt.Errorf("request in error state")

	rn = []byte("\r\n")
)

type parserState string

const (
	StateInit    parserState = "init"
	StateHeaders parserState = "headers"
	StateDone    parserState = "done"
	StateError   parserState = "error"
)

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

type Request struct {
	RequestLine RequestLine
	state       parserState
	Headers     *headers.Headers
}

func newRequest() *Request {
	return &Request{
		state:   StateInit,
		Headers: headers.NewHeaders(),
	}
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, rn)
	if idx == -1 {
		return nil, 0, nil
	}
	startLine := b[:idx]
	read := idx + len(rn)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ERROR_BAD_START_LINE
	}

	splits := bytes.Split(parts[2], []byte("/"))
	if len(splits) != 2 || string(splits[0]) != "HTTP" || string(splits[1]) != "1.1" {
		return nil, 0, ERROR_BAD_START_LINE
	}
	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(splits[1]),
	}
	return rl, read, nil
}

func (req *Request) parse(data []byte) (int, error) {
	var read int
outer:
	for {
		currentData := data[read:]
		switch req.state {

		case StateInit:
			rl, n, err := parseRequestLine(currentData)
			if err != nil {
				req.state = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			req.RequestLine = *rl
			read += n
			req.state = StateHeaders

		case StateHeaders:
			n, done, err := req.Headers.Parse(currentData)
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break outer
			}
			read += n
			if done {
				req.state = StateDone
			}

		case StateDone:
			break outer

		case StateError:
			return 0, ERROR_REQUEST_IN_ERROR_STATE
		default:
			panic("somehow we have programmed poorly")
		}
	}
	return read, nil
}

func (req *Request) done() bool {
	return req.state == StateDone || req.state == StateError
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := newRequest()

	buf := make([]byte, 1024)
	var bufLen int

	for !req.done() {
		n1, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}
		bufLen += n1
		n2, err := req.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[n2:bufLen])
		bufLen -= n2
	}
	return req, nil
}
