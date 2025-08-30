package headers

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/kalogs-c/the-go-http/internal/constraints"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(map[string]string)
}

var (
	ErrorExtraSpaceBeforeColon = errors.New("extra space before colon")
	ErrorInvalidCharOnFieldKey = errors.New("invalid char on field")
)

var allowedCharsSet map[rune]bool = craftAllowedCharSet()

func craftAllowedCharSet() map[rune]bool {
	allowedChars := "!#$%&'*+-.^_`|~"
	for c := '0'; c <= '9'; c++ {
		allowedChars += string(c)
	}

	for c := 'a'; c <= 'z'; c++ {
		allowedChars += string(c)
	}

	for c := 'A'; c <= 'Z'; c++ {
		allowedChars += string(c)
	}

	allowedCharsSet := make(map[rune]bool)
	for _, c := range allowedChars {
		allowedCharsSet[c] = true
	}

	return allowedCharsSet
}

func validateHeaderField(field string) error {
	if field[len(field)-1] == ' ' {
		return ErrorExtraSpaceBeforeColon
	}

	i := strings.IndexFunc(field, func(c rune) bool {
		_, ok := allowedCharsSet[c]
		return !ok
	})

	if i != -1 {
		errorMsg := fmt.Sprintf("invalid char '%c' at index %d: %s\n", field[i], i, field)
		errorMsg += fmt.Sprintf("%*s", len(errorMsg)+i-len(field), "^")
		errorMsg += " invalid char here"

		return errors.Join(ErrorInvalidCharOnFieldKey, errors.New(errorMsg))
	}

	return nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	i := bytes.Index(data, constraints.CRLF)
	switch i {
	case 0: // End of headers
		return len(constraints.CRLF), true, nil
	case -1: // Incomplete headers
		return 0, false, nil
	}

	rawContent := strings.TrimSpace(string(data[:i]))
	parts := strings.SplitN(rawContent, ":", 2)

	if err := validateHeaderField(parts[0]); err != nil {
		return 0, false, err
	}

	field := strings.TrimSpace(strings.ToLower(parts[0]))
	value := strings.TrimSpace(parts[1])

	if _, ok := h[field]; ok {
		h[field] += ", " + value
	} else {
		h[field] = value
	}

	return i + len(constraints.CRLF), false, nil
}
