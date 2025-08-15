// Package request
package request

import (
	"bytes"
	"errors"
	"io"
	"slices"
)

type RequestLine struct {
	HTTPVersion   string
	RequestTarget string
	Method        string
}

type parserState string

const (
	StateInit parserState = "init"
	StateDone parserState = "done"
)

type Request struct {
	RequestLine RequestLine
	state       parserState
}

func (r *Request) done() bool {
	return r.state == StateDone
}

func newRequest() *Request {
	return &Request{state: StateInit}
}

func (r *Request) parse(data []byte) (int, error) {
	if r.state == StateDone {
		return 0, nil
	}

	rl, n, err := parseRequestLine(data)
	if err != nil {
		return n, err
	}

	if rl == nil {
		return 0, nil
	}

	r.RequestLine = *rl
	r.state = StateDone

	return n, nil
}

var CRLF = []byte("\r\n")
var (
	ErrMalformedRequestLine   = errors.New("malformed request-line")
	ErrUnsupportedHTTPVersion = errors.New("unsupported http version")
	ErrInvalidMethod          = errors.New("invalid method")
)

var methods = []string{"GET", "POST"}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, CRLF)
	if idx == -1 {
		return nil, 0, nil
	}

	requestLine := b[:idx]
	read := len(requestLine) + len(CRLF)

	parts := bytes.Split(requestLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ErrMalformedRequestLine
	}

	if !slices.Contains(methods, string(parts[0])) {
		return nil, 0, ErrInvalidMethod
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" {
		return nil, 0, ErrMalformedRequestLine
	}
	if string(httpParts[1]) != "1.1" {
		return nil, 0, ErrUnsupportedHTTPVersion
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HTTPVersion:   string(httpParts[1]),
	}

	return rl, read, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}
		bufLen += n

		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	return request, nil
}
