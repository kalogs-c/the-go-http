package request

import (
	"errors"
	"strings"
	"unsafe"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var (
	ERROR_REQUEST_LINE_INCOMPLETE  error = errors.New("number of parts does not match, it should be 'METHOD REQUEST-TARGET HTTP-VERSION'")
	ERROR_METHOD_NOT_ALLOWED       error = errors.New("the method should only contains capital alphabetic characters")
	ERROR_HTTP_VERSION_NOT_ALLOWED error = errors.New("http version is not suported")
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
	i := strings.Index(content, separator)
	if i == -1 {
		return nil, 0, "", nil
	}

	startLine := content[:i]
	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, 0, "", ERROR_REQUEST_LINE_INCOMPLETE
	}

	versionParts := strings.Split(parts[2], "/")

	reqLine := RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   versionParts[1],
	}

	if !isMethodAllowed(reqLine.Method) {
		return nil, 0, "", ERROR_METHOD_NOT_ALLOWED
	}

	if !isHttpVersionAllowed(reqLine.HttpVersion) {
		return nil, 0, "", ERROR_HTTP_VERSION_NOT_ALLOWED
	}

	restOfMsg := content[i+len(separator):]
	return &reqLine, int(unsafe.Sizeof(startLine)), restOfMsg, nil
}
