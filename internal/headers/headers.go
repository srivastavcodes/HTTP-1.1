package headers

import (
	"bytes"
	"fmt"
)

var rn = []byte("\r\n")

type Headers map[string]string

func NewHeaders() Headers {
	return make(map[string]string)
}

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

func (hdr Headers) Parse(data []byte) (int, bool, error) {
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
		read += idx + len(rn)
		hdr[name] = value
	}
	return read, done, nil
}
