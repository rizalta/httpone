// Package headers
package headers

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"unicode"
)

type Headers map[string]string

func (h Headers) Get(name string) string {
	return h[strings.ToLower(name)]
}

func (h Headers) Set(name, value string) {
	name = strings.ToLower(name)
	if v, ok := h[name]; ok {
		value = fmt.Sprintf("%s, %s", v, value)
	}
	h[name] = value
}

func NewHeaders() Headers {
	return make(Headers)
}

var (
	ErrMalformedHeader   = errors.New("malformed header")
	ErrInvalidHeaderName = errors.New("invalid header name")
)

func isToken(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			continue
		}

		switch r {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			continue
		default:
			return false
		}
	}
	return true
}

var crlf = []byte("\r\n")

func (h Headers) Parse(data []byte) (int, bool, error) {
	n := 0
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
		return 0, false, ErrMalformedHeader
	}

	name := bytes.ToLower(bytes.TrimSpace(line[:colonIdx]))

	if !isToken(string(name)) {
		return 0, false, ErrInvalidHeaderName
	}

	value := bytes.TrimSpace(line[colonIdx+1:])

	h.Set(string(name), string(value))
	return n, false, nil
}
