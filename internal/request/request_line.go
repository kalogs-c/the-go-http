package request

import (
	"errors"
	"strings"
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

func isMethodAllowed(method string) bool {
	return !strings.ContainsFunc(method, func(c rune) bool {
		return (c < 'A' || c > 'Z')
	})
}

func isHttpVersionAllowed(httpVersion string) bool {
	return httpVersion == "1.1"
}

func parseRequestLine(content string) (*RequestLine, int, string, error) {
	i := strings.Index(content, crlf)
	if i == -1 {
		return nil, 0, "", nil
	}

	startLine := content[:i]
	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, 0, "", ErrorRequestLineIncomplete
	}

	versionParts := strings.Split(parts[2], "/")

	reqLine := RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   versionParts[1],
	}

	if !isMethodAllowed(reqLine.Method) {
		return nil, 0, "", ErrorMethodNotAllowed
	}

	if !isHttpVersionAllowed(reqLine.HttpVersion) {
		return nil, 0, "", ErrorHttpVersionNotAllowed
	}

	read := i + len(crlf)
	restOfMsg := content[read:]
	return &reqLine, read, restOfMsg, nil
}
