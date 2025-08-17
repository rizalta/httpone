package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseHeaders(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Valid header with extra spacing
	headers = NewHeaders()
	data = []byte("          Host: localhost:42069    \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers.Get("host"))
	assert.Equal(t, 37, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	headers = NewHeaders()
	headers.Set("host", "localhost:42069")
	data = []byte("User-Agent: curl\r\nConnection: Keep-Me\r\n\r\n")
	read := 0
	for {
		n, done, err := headers.Parse(data[read:])
		require.NoError(t, err)
		read += n
		if done {
			break
		}
	}
	assert.Equal(t, "localhost:42069", headers.Get("host"))
	assert.Equal(t, "curl", headers.Get("user-agent"))
	assert.Equal(t, "Keep-Me", headers.Get("connection"))
	assert.Equal(t, 41, read)

	// Test: Invalid character in header
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.False(t, done)
	assert.Equal(t, 0, n)

	// Test: Valid headers with digits and special charcters
	headers = NewHeaders()
	data = []byte("    X-Token_123!#$%&'*+.^`|~:       Testing     \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.False(t, done)
	assert.Equal(t, 50, n)
	assert.Equal(t, "Testing", headers.Get("x-token_123!#$%&'*+.^`|~"))

	// Test: Multiple values for the same name
	headers = NewHeaders()
	data = []byte("Foo: Bar\r\nfoo: notBar\r\n\r\n")
	read = 0
	for {
		n, done, err := headers.Parse(data[read:])
		require.NoError(t, err)
		read += n
		if done {
			break
		}
	}
	assert.Equal(t, "Bar, notBar", headers.Get("foo"))
	assert.Equal(t, 25, read)
}
