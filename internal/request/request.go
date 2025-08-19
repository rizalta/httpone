// Package request
package request

import (
	"bytes"
	"errors"
	"io"
	"strconv"

	"github.com/rizalta/httpone/internal/headers"
)

type RequestLine struct {
	HTTPVersion   string
	RequestTarget string
	Method        string
}

type parserState string

const (
	StateInit    parserState = "init"
	StateHeaders parserState = "headers"
	StateBody    parserState = "body"
	StateDone    parserState = "done"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	state       parserState
}

func (r *Request) done() bool {
	return r.state == StateDone
}

func (r *Request) contentLength() int {
	contentLenStr := r.Headers.Get("content-length")
	if contentLenStr == "" {
		return 0
	}

	contentLen, err := strconv.Atoi(contentLenStr)
	if err != nil {
		return 0
	}
	return max(0, contentLen)
}

func newRequest() *Request {
	return &Request{
		state:   StateInit,
		Headers: headers.NewHeaders(),
	}
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case StateDone:
		return 0, nil

	case StateInit:
		rl, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}

		r.RequestLine = *rl
		r.state = StateHeaders

		return n, nil

	case StateHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			if r.contentLength() > 0 {
				r.state = StateBody
			} else {
				r.state = StateDone
			}
		}
		return n, nil

	case StateBody:
		length := r.contentLength()
		if length == 0 {
			r.state = StateDone
			return 0, nil
		}

		needed := length - len(r.Body)
		available := len(data)
		read := min(needed, available)

		r.Body = append(r.Body, data[:read]...)
		if len(r.Body) == length {
			r.state = StateDone
		}

		return read, nil
	}

	return 0, nil
}

var crlf = []byte("\r\n")
var (
	ErrMalformedRequestLine   = errors.New("malformed request-line")
	ErrUnsupportedHTTPVersion = errors.New("unsupported http version")
	ErrInvalidMethod          = errors.New("invalid method")
)

var methods = map[string]struct{}{
	"GET":     {},
	"POST":    {},
	"PUT":     {},
	"DELETE":  {},
	"HEAD":    {},
	"OPTIONS": {},
	"PATCH":   {},
}

func isValidMethod(m string) bool {
	_, ok := methods[m]
	return ok
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, crlf)
	if idx == -1 {
		return nil, 0, nil
	}

	requestLine := b[:idx]
	read := len(requestLine) + len(crlf)

	parts := bytes.Split(requestLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ErrMalformedRequestLine
	}

	if !isValidMethod(string(parts[0])) {
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

	buf := &bytes.Buffer{}

	for !request.done() {
		data := buf.Bytes()
		readN, err := request.parse(data)
		if err != nil {
			return nil, err
		}

		if readN > 0 {
			buf.Next(readN)
			continue
		}

		chunk := make([]byte, 1024)
		n, err := reader.Read(chunk)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}

		if n > 0 {
			buf.Write(chunk[:n])
		}

		if errors.Is(err, io.EOF) && n == 0 {
			if buf.Len() > 0 {
				_, parseErr := request.parse(buf.Bytes())
				if parseErr != nil {
					return nil, parseErr
				}
			}

			if !request.done() {
				return nil, io.ErrUnexpectedEOF
			}
			break
		}
	}

	return request, nil
}
