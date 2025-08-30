package main

import (
	"fmt"
	"log"
	"net"

	"github.com/kalogs-c/the-go-http/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("Request line:")
		fmt.Printf(" - Method:  %s\n", r.RequestLine.Method)
		fmt.Printf(" - Target:  %s\n", r.RequestLine.RequestTarget)
		fmt.Printf(" - Version: %s\n", r.RequestLine.HttpVersion)

		fmt.Println("Headers:")
		for key, value := range r.Headers {
			fmt.Printf(" - %s: %s\n", key, value)
		}

		conn.Close()
	}
}
