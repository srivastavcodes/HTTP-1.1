package headers

import (
	"bytes"
	"fmt"
	"strings"
)

func isToken(str []byte) bool {
	for _, ch := range str {
		var is bool

		if ch > 'A' && ch < 'Z' || ch > 'a' && ch < 'z' || ch > '0' && ch < '9' {
			is = true
		}
		switch ch {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			is = true
		}
		if !is {
			return false
		}
	}
	return true
}

var rn = []byte("\r\n")

func parseHeader(line []byte) (string, string, error) {
	parts := bytes.SplitN(line, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed field line")
	}
	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("malformed field name")
	}
	return string(name), string(value), nil
}

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: make(map[string]string),
	}
}

func (hdr *Headers) Get(name string) string {
	return hdr.headers[strings.ToLower(name)]
}

func (hdr *Headers) Set(name, value string) {
	hdr.headers[strings.ToLower(name)] = value
}

func (hdr *Headers) Parse(data []byte) (int, bool, error) {
	done, read := false, 0
	for {
		idx := bytes.Index(data[read:], rn)
		if idx == -1 {
			break
		}
		if idx == 0 {
			read += len(rn)
			done = true
			break
		}
		name, value, err := parseHeader(data[read : read+idx])
		if err != nil {
			return 0, false, err
		}
		if !isToken([]byte(name)) {
			return 0, false, fmt.Errorf("malformed header name")
		}
		read += idx + len(rn)
		hdr.Set(name, value)
	}
	return read, done, nil
}
