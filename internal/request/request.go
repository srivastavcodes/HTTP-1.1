package request

import (
	"bytes"
	"fmt"
	"io"
)

var (
	ERROR_BAD_START_LINE         = fmt.Errorf("malformed request-line")
	ERROR_REQUEST_IN_ERROR_STATE = fmt.Errorf("request in error state")

	SEPARATOR = "\r\n"
)

type parserState string

const (
	StateInit  parserState = "init"
	StateDone  parserState = "done"
	StateError parserState = "error"
)

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

type Request struct {
	RequestLine RequestLine
	state       parserState
}

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, []byte(SEPARATOR))
	if idx == -1 {
		return nil, 0, nil
	}
	startLine := b[:idx]
	read := idx + len(SEPARATOR)

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
		switch req.state {
		case StateError:
			return 0, ERROR_REQUEST_IN_ERROR_STATE

		case StateInit:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				req.state = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			req.RequestLine = *rl
			read += n
			req.state = StateDone

		case StateDone:
			break outer
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
