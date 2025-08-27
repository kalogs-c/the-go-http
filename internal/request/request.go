package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const SEPARATOR = "\r\n"

func parseRequestLine(content string) (*RequestLine, string, error) {
	i := strings.Index(content, SEPARATOR)
	if i == -1 {
		return nil, "", errors.New("request line not found")
	}

	startLine := content[:i]
	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, "", errors.New("number of parts does not match, it should be 'METHOD REQUEST-TARGET HTTP-VERSION'")
	}

	versionParts := strings.Split(parts[2], "/")

	reqLine := RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   versionParts[1],
	}

	restOfMsg := content[i+len(SEPARATOR):]
	return &reqLine, restOfMsg, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	reqLine, _, err := parseRequestLine(string(content))
	if err != nil {
		return nil, err
	}

	req := Request{
		RequestLine: *reqLine,
	}

	return &req, nil
}
