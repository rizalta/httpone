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

type response struct {
	wroteHeader bool
	headers     headers.Headers
	writer      io.Writer
}

func NewResponse(w io.Writer) *response {
	return &response{
		writer:  w,
		headers: GetDefaultHeaders(),
	}
}

func (w *response) WriteHeader(statusCode StatusCode) error {
	w.wroteHeader = true
	header := fmt.Appendf(nil, "HTTP/1.1 %d %s\r\n", statusCode, statusMessage[statusCode])
	for n := range w.headers {
		header = fmt.Appendf(header, "%s: %s\r\n", formatHeeaderName(n), w.headers.Get(n))
	}
	_, err := w.writer.Write(header)
	return err
}

func (w *response) Write(p []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(StatusOK)
	}
	contentLenHeader := fmt.Appendf(nil, "%s: %d\r\n", "Content-Length", len(p))
	w.writer.Write(contentLenHeader)
	w.writer.Write([]byte("\r\n"))
	return w.writer.Write(p)
}

func GetDefaultHeaders() headers.Headers {
	h := headers.NewHeaders()
	h.Set("connection", "close")
	h.Set("content-type", "text/plain")

	return h
}

func formatHeeaderName(name string) string {
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
