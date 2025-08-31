package request

import (
	"bytes"
	"errors"
	"strings"

	"github.com/kalogs-c/the-go-http/internal/constraints"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var (
	ErrorRequestLineIncomplete error = errors.New("number of parts does not match, it should be 'METHOD REQUEST-TARGET HTTP-VERSION'")
	ErrorMethodNotAllowed      error = errors.New("the method should only contains capital alphabetic characters")
	ErrorHttpVersionNotAllowed error = errors.New("http version is not suported")
)

func methodAllowed(method string) bool {
	return !strings.ContainsFunc(method, func(c rune) bool {
		return (c < 'A' || c > 'Z')
	})
}

func httpVersionAllowed(httpVersion string) bool {
	return httpVersion == "1.1"
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	i := bytes.Index(data, constraints.CRLF)
	if i == -1 {
		return nil, 0, nil
	}

	startLine := string(data[:i])
	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, 0, ErrorRequestLineIncomplete
	}

	versionParts := strings.Split(parts[2], "/")

	requestLine := RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   versionParts[1],
	}

	if !methodAllowed(requestLine.Method) {
		return nil, 0, ErrorMethodNotAllowed
	}

	if !httpVersionAllowed(requestLine.HttpVersion) {
		return nil, 0, ErrorHttpVersionNotAllowed
	}

	bytesRead := i + len(constraints.CRLF)
	return &requestLine, bytesRead, nil
}
