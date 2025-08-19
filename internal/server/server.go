// Package server
package server

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/rizalta/httpone/internal/request"
	"github.com/rizalta/httpone/internal/response"
)

type Server struct {
	listener net.Listener
	handler  Handler
	closed   atomic.Bool
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: listener,
		handler:  handler,
	}

	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	if s.closed.CompareAndSwap(false, true) {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				log.Println("server stopped")
				return
			}
			log.Printf("error accepting conn, %v\n", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	fmt.Printf("req: %v\n", req)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.StatusInternalServerError,
			Message:    err.Error(),
		}
		hErr.Write(conn)
		return
	}
	buf := bytes.NewBuffer([]byte{})
	hErr := s.handler(buf, req)
	if hErr != nil {
		hErr.Write(conn)
		return
	}
	b := buf.Bytes()
	headers := response.GetDefaultHeaders(len(b))
	response.WriteStatusLine(conn, response.StatusOK)
	response.WriteHeaders(conn, headers)
	conn.Write(b)
}
