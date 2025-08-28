package request

import (
	"errors"
	"io"
	"strings"
	"unsafe"
)

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateDone
)

type Request struct {
	RequestLine RequestLine
	state       requestState
}

func newRequest() *Request {
	return &Request{state: requestStateInitialized}
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const (
	separator  string = "\r\n"
	bufferSize int    = 8
)

var (
	ERROR_REQUEST_LINE_INCOMPLETE  error = errors.New("number of parts does not match, it should be 'METHOD REQUEST-TARGET HTTP-VERSION'")
	ERROR_METHOD_NOT_ALLOWED       error = errors.New("the method should only contains capital alphabetic characters")
	ERROR_HTTP_VERSION_NOT_ALLOWED error = errors.New("http version is not suported")

	ERROR_READ_ON_DONE_STATE error = errors.New("trying to read data in a done state")
)

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		reqLine, bytesRead, _, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}

		if bytesRead == 0 {
			return 0, nil
		}

		r.RequestLine = *reqLine
		r.state = requestStateDone
		return bytesRead, nil
	case requestStateDone:
		return 0, ERROR_READ_ON_DONE_STATE
	default:
		return 0, errors.New("unknown state")
	}
}

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

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, 0, bufferSize)
	readToIndex := 0
	request := newRequest()

	for request.state != requestStateDone {
		tmpBuf := make([]byte, bufferSize)
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
