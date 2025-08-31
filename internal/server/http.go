package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"
	"time"

	"github.com/kalogs-c/the-go-http/internal/request"
	"github.com/kalogs-c/the-go-http/internal/response"
)

type Server struct {
	listener net.Listener
	closed   *atomic.Bool
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	var closed atomic.Bool
	closed.Store(false)
	server := Server{listener, &closed, handler}

	go server.listen()

	return &server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for !s.closed.Load() {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		conn.SetReadDeadline(time.Now().Add(time.Second * 5))

		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("request from %s to %s accepted", conn.RemoteAddr().String(), r.RequestLine.RequestTarget)
		handlerErr := s.handle(conn, r)
		if handlerErr != nil {
			log.Fatalf(
				"request to '%s' failed with status code %d: %s",
				r.RequestLine.RequestTarget,
				handlerErr.StatusCode,
				handlerErr.Message,
			)
		}
	}
}

func (s *Server) handle(conn net.Conn, r *request.Request) *HandlerError {
	h := response.GetDefaultHeaders(13)

	io.WriteString(conn, "HTTP/1.1 200 OK\r\n")
	response.WriteHeaders(conn, h)
	err := s.handler(conn, r)
	conn.Close()
	return err
}
