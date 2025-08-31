package server

import (
	"io"

	"github.com/kalogs-c/the-go-http/internal/request"
	"github.com/kalogs-c/the-go-http/internal/response"
)

type (
	Handler      func(w io.Writer, req *request.Request) *HandlerError
	HandlerError struct {
		StatusCode response.StatusCode
		Message    string
	}
)
