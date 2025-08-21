// Package response
package response

import (
	"fmt"
	"io"
	"unicode"

	"github.com/rizalta/httpone/internal/headers"
)

type Writer interface {
	WriteHeader(statusCode StatusCode) error
	Write([]byte) (int, error)
	Headers() headers.Headers
}

type writerState int

const (
	stateInit writerState = iota
	stateHeader
	stateBody
	stateDone
)

type response struct {
	state   writerState
	headers headers.Headers
	writer  io.Writer
}

func NewResponse(w io.Writer) *response {
	return &response{
		writer:  w,
		headers: GetDefaultHeaders(),
		state:   stateInit,
	}
}

func (w *response) WriteHeader(statusCode StatusCode) error {
	header := fmt.Appendf(nil, "HTTP/1.1 %d %s\r\n", statusCode, statusMessage[statusCode])
	for n := range w.headers {
		header = fmt.Appendf(header, "%s: %s\r\n", formatHeaderName(n), w.headers.Get(n))
	}
	_, err := w.writer.Write(header)
	w.state = stateHeader
	return err
}

func (w *response) Write(p []byte) (int, error) {
	if w.state == stateInit {
		w.WriteHeader(StatusOK)
	}
	contentLenHeader := fmt.Appendf(nil, "%s: %d\r\n\r\n", "Content-Length", len(p))
	_, err := w.writer.Write(contentLenHeader)
	if err != nil {
		return 0, err
	}
	w.state = stateBody
	return w.writer.Write(p)
}

func GetDefaultHeaders() headers.Headers {
	h := headers.NewHeaders()
	h.Set("connection", "close")
	h.Set("content-type", "text/plain")

	return h
}

func formatHeaderName(name string) string {
	runes := []rune(name)
	delimiter := '-'
	capNext := true
	for i, r := range runes {
		if r == delimiter {
			capNext = true
			continue
		}
		if capNext {
			runes[i] = unicode.ToUpper(r)
			capNext = false
		}
	}

	return string(runes)
}

func (w *response) Headers() headers.Headers {
	return w.headers
}

func (w *response) Flush() {
	switch w.state {
	case stateBody:
		return
	case stateHeader:
		w.writer.Write([]byte("\r\n"))
	case stateInit:
		w.WriteHeader(StatusOK)
		w.writer.Write([]byte("\r\n"))
	}
}
