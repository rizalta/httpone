package main

import (
	"fmt"
	"log"
	"net"

	"github.com/rizalta/httpone/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("connection not accepted")
			break
		}
		defer conn.Close()
		log.Println("connection accepted")

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Printf("error reading request %v\n", err)
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n\n", req.RequestLine.HTTPVersion)

		log.Println("connection closed")
		conn.Close()
	}
}
