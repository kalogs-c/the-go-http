package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(map[string]string)
}

const crlf string = "\r\n"

var ErrorExtraSpaceBeforeColon = errors.New("extra space before colon")

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	i := strings.Index(string(data), crlf)
	switch i {
	case 0: // End of headers
		return len(crlf), true, nil
	case -1: // Incomplete headers
		return 0, false, nil
	}

	rawContent := strings.TrimSpace(string(data[:i]))
	parts := strings.SplitN(rawContent, ":", 2)

	field := parts[0]
	value := strings.TrimSpace(parts[1])

	if field[len(field)-1] == ' ' {
		return 0, false, ErrorExtraSpaceBeforeColon
	}

	field = strings.TrimSpace(field)

	h[field] = value

	return i + len(crlf), false, nil
}
