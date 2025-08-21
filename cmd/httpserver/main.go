package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rizalta/httpone/internal/request"
	"github.com/rizalta/httpone/internal/response"
	"github.com/rizalta/httpone/internal/server"
)

const port = 42069

var html400 = []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
	</html>`)

var html500 = []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)

var html200 = []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)

func main() {
	server, err := server.Serve(port, func(w response.Writer, req *request.Request) {
		var status response.StatusCode = response.StatusOK
		var body []byte
		w.Headers().Set("Content-Type", "text/html")
		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			status = response.StatusBadRequest
			body = html400
			w.WriteHeader(status)
			w.Write(body)
		case "/myproblem":
			status = response.StatusInternalServerError
			body = html500
			w.WriteHeader(status)
			w.Write(body)
		}
	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
