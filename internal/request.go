package internal

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	ERROR_BAD_START_LINE = fmt.Errorf("malformed request-line")
)

var SEPARATOR = "\r\n"

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

type Request struct {
	RequestLine RequestLine
}

func parseRequestLine(str string) (*RequestLine, string, error) {
	idx := strings.Index(str, SEPARATOR)
	if idx == -1 {
		return nil, str, nil
	}
	startLine := str[:idx]
	remaining := str[idx+len(SEPARATOR):]

	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, remaining, ERROR_BAD_START_LINE
	}
	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return nil, remaining, ERROR_BAD_START_LINE
	}
	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   httpParts[1],
	}
	return rl, remaining, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("unable to io.ReadAll"), err)
	}
	str := string(data)

	rl, _, err := parseRequestLine(str)
	if err != nil {
		return nil, err
	}
	return &Request{RequestLine: *rl}, nil
}
