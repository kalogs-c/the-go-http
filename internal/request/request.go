package request

import (
	"errors"
	"io"

	"github.com/kalogs-c/the-go-http/internal/constraints"
	"github.com/kalogs-c/the-go-http/internal/headers"
)

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateParsingHeaders
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

		r.state = requestStateDone
		return bytesRead, nil
	case requestStateDone:
		return 0, ErrorReadOnDoneState
	default:
		return 0, errors.New("unknown state")
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, 0, constraints.BufferSize)
	readToIndex := 0
	request := newRequest()

	for request.state != requestStateDone {
		tmpBuf := make([]byte, constraints.BufferSize)
		n, readErr := reader.Read(tmpBuf)
		if readErr != nil && readErr != io.EOF {
			return nil, readErr
		}

		buffer = append(buffer, tmpBuf[:n]...)

		readToIndex = n
		n, err := request.parse(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}

		if readErr == io.EOF {
			request.state = requestStateDone
		}

		if n > 0 {
			buffer = buffer[n:]
			readToIndex -= n
		}
	}

	return request, nil
}
