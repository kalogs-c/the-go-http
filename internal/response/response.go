package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/kalogs-c/the-go-http/internal/constraints"
	"github.com/kalogs-c/the-go-http/internal/headers"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusNotFound            StatusCode = 404
	StatusInternalServerError StatusCode = 500
)

func GetDefaultHeaders(contentLength int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", strconv.Itoa(contentLength))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

func WriteHeaders(w io.Writer, h headers.Headers) error {
	for key, value := range h {
		header := fmt.Sprintf("%s: %s%s", key, value, constraints.CRLF)
		_, err := io.WriteString(w, header)
		if err != nil {
			return err
		}
	}

	_, err := w.Write(constraints.CRLF)
	return err
}
