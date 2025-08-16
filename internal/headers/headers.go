// Package headers
package headers

import (
	"bytes"
	"errors"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

var ErrInvalidFieldLine = errors.New("invalid field-line")

var crlf = []byte("\r\n")

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, crlf)
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return len(crlf), true, nil
	}
	line := data[:idx]
	n += idx + len(crlf)
	colonIdx := bytes.IndexByte(line, ':')
	if colonIdx <= 0 || line[colonIdx-1] == ' ' {
		return 0, false, ErrInvalidFieldLine
	}
	key := bytes.TrimSpace(line[:colonIdx])
	value := bytes.TrimSpace(line[colonIdx+1:])
	h[string(key)] = string(value)
	return n, false, nil
}
