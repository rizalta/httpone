package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", err)
	}

	for {
		conn, err := listener.Accept()
		log.Println("connection accepted")

		if err != nil {
			log.Fatal("error", err)
		}

		for line := range getLinesChannel(conn) {
			fmt.Printf("%s\n", line)
		}
		log.Println("connection closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(ch)

		s := ""
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				break
			}

			data = data[:n]
			if i := strings.IndexByte(string(data), '\n'); i != -1 {
				s += string(data[:i])
				data = data[i+1:]
				ch <- s
				s = ""
			}

			s += string(data)
		}

		if len(s) != 0 {
			ch <- s
		}
	}()

	return ch
}
