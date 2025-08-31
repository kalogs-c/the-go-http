package request

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/kalogs-c/the-go-http/internal/constraints"
	"github.com/kalogs-c/the-go-http/internal/headers"
)

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	state       requestState
}

func newRequest() *Request {
	return &Request{
		Headers: headers.NewHeaders(),
		state:   requestStateInitialized,
	}
}

var ErrorReadOnDoneState error = errors.New("trying to read data in a done state")

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		reqLine, bytesRead, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if bytesRead == 0 {
			return 0, nil
		}

		r.RequestLine = *reqLine
		r.state = requestStateParsingHeaders
		return bytesRead, nil
	case requestStateParsingHeaders:
		bytesRead, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}

		if !done {
			return bytesRead, nil
		}

		r.state = requestStateParsingBody
		return bytesRead, nil
	case requestStateParsingBody:
		contentLengthStr, ok := r.Headers.Get("Content-Length")
		if !ok {
			r.state = requestStateDone
			return 0, nil
		}

		contentLength, err := strconv.Atoi(contentLengthStr)
		if err != nil {
			return 0, err
		}

		r.Body = append(r.Body, data...)
		if len(r.Body) > contentLength {
			return 0, errors.New("body bigger than Content-Length passed")
		}

		if len(r.Body) == contentLength {
			r.state = requestStateDone
		}

		if len(data) == 0 && len(r.Body) < contentLength {
			return 0, fmt.Errorf("body length mismatch: expected %d, got %d", contentLength, len(r.Body))
		}

		return len(data), nil
	case requestStateDone:
		return 0, ErrorReadOnDoneState
	default:
		return 0, errors.New("unknown state")
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, 0, constraints.BufferSize)
	request := newRequest()

	for request.state != requestStateDone {
		tmpBuf := make([]byte, constraints.BufferSize)
		n, readErr := reader.Read(tmpBuf)
		if readErr != nil && readErr != io.EOF {
			return nil, readErr
		}

		buffer = append(buffer, tmpBuf[:n]...)

		n, err := request.parse(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}

		if readErr == io.EOF {
			request.state = requestStateDone
		}

		if n > 0 {
			buffer = buffer[n:]
		}
	}

	return request, nil
}
