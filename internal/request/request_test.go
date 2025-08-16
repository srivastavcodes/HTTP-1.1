package request

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// It's useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n
	if n > cr.numBytesPerRead {
		n = cr.numBytesPerRead
		cr.pos -= n - cr.numBytesPerRead
	}
	return n, nil
}

func TestRequestLineParse(t *testing.T) {
	// Test: Good GET Request line
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}
	req, err := RequestFromReader(reader)

	require.NoError(t, err)
	require.NotNil(t, req)

	assert.Equal(t, "GET", req.RequestLine.Method)
	assert.Equal(t, "/", req.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", req.RequestLine.HttpVersion)

	// Test: Good GET Request line with path
	reader = &chunkReader{
		data:            "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 1,
	}
	req, err = RequestFromReader(reader)

	require.NoError(t, err)
	require.NotNil(t, req)

	assert.Equal(t, "GET", req.RequestLine.Method)
	assert.Equal(t, "/coffee", req.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", req.RequestLine.HttpVersion)
}

func TestParseHeaders(t *testing.T) {
	// Test: Standard Headers
	reader := &chunkReader{numBytesPerRead: 3,
		data: "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	}
	req, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, req)

	host, ok := req.Headers.Get("host")
	assert.True(t, ok)
	assert.Equal(t, "localhost:42069", host)

	userAgent, ok := req.Headers.Get("user-agent")
	assert.True(t, ok)
	assert.Equal(t, "curl/7.81.0", userAgent)

	accept, ok := req.Headers.Get("accept")
	assert.True(t, ok)
	assert.Equal(t, "*/*", accept)

	// Test: Malformed Header
	reader = &chunkReader{numBytesPerRead: 3,
		data: "GET / HTTP/1.1\r\nHost localhost:42069\r\n\r\n",
	}
	req, err = RequestFromReader(reader)
	require.Error(t, err)
}

func TestParseBody(t *testing.T) {
	// Test: Standard Body
	reader := &chunkReader{
		data: "POST /submit HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"Content-Length: 13\r\n" +
			"\r\n" +
			"hello world!\n",
		numBytesPerRead: 3,
	}
	req, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, req)
	assert.Equal(t, "hello world!\n", string(req.Body))

	// Test: Body shorter than reported content length
	reader = &chunkReader{
		data: "POST /submit HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"Content-Length: 20\r\n" +
			"\r\n" +
			"partial content",
		numBytesPerRead: 3,
	}
	req, err = RequestFromReader(reader)
	require.Error(t, err)
}
