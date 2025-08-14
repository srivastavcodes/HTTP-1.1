package internal

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestLineParse(t *testing.T) {
	// Test: Good GET Request line
	reader := strings.NewReader("GET / HTTP/1.1\r\n" +
		"Host: localhost:7714\r\nUser-Agent: curl/8.7.1\r\nAccept: */*\r\n\r\n")

	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)

	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Good GET Request line with path
	reader = strings.NewReader("GET /coffee HTTP/1.1\r\n" +
		"Host: localhost:7714\r\nUser-Agent: curl/8.7.1\r\nAccept: */*\r\n\r\n")

	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)

	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Invalid number of parts in request line
	reader = strings.NewReader("/coffee HTTP/1.1\r\n" +
		"Host: localhost:7714\r\nUser-Agent: curl/8.7.1\r\nAccept: */*\r\n\r\n")

	_, err = RequestFromReader(reader)
	require.Error(t, err)
}
