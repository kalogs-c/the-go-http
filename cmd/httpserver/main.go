package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kalogs-c/the-go-http/internal/request"
	"github.com/kalogs-c/the-go-http/internal/server"
)

const port = 42069

func handleRoot(w io.Writer, r *request.Request) *server.HandlerError {
	io.WriteString(w, "hello!\n")
	return nil
}

func main() {
	server, err := server.Serve(port, handleRoot)
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
