// Package response
package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/rizalta/httpone/internal/headers"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := fmt.Appendf(nil, "HTTP/1.1 %d ", statusCode)
	switch statusCode {
	case StatusOK:
		statusLine = fmt.Appendf(statusLine, "OK\r\n")
	case StatusBadRequest:
		statusLine = fmt.Appendf(statusLine, "Bad Request\r\n")
	case StatusInternalServerError:
		statusLine = fmt.Appendf(statusLine, "Internal Server Error\r\n")
	}

	_, err := w.Write(statusLine)
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()

	h.Set("content-length", strconv.Itoa(contentLen))
	h.Set("content-type", "text/plain")

	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	b := []byte{}
	for n := range headers {
		b = fmt.Appendf(b, "%s: %s\r\n", n, headers.Get(n))
	}
	b = fmt.Appendf(b, "\r\n")

	_, err := w.Write(b)

	return err
}
