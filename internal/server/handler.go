package server

import (
	"io"

	"github.com/rizalta/httpone/internal/request"
	"github.com/rizalta/httpone/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (he *HandlerError) Write(w io.Writer) error {
	if err := response.WriteStatusLine(w, he.StatusCode); err != nil {
		return err
	}

	defaultHeaders := response.GetDefaultHeaders(len(he.Message))
	if err := response.WriteHeaders(w, defaultHeaders); err != nil {
		return err
	}

	_, err := w.Write([]byte(he.Message))
	return err
}
