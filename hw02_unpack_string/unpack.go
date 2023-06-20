package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

const (
	zero   int  = int(rune('0'))
	escape rune = '\\'
	none   rune = 0
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	sb := strings.Builder{}

	var symbol rune
	escaping := false
	for _, currentSymbol := range input {
		switch {
		case unicode.IsDigit(currentSymbol):
			if escaping {
				symbol = currentSymbol
				escaping = false
				continue
			}

			if symbol == none {
				return "", ErrInvalidString
			}

			sb.WriteString(strings.Repeat(string(symbol), int(currentSymbol)-zero))
			symbol = none

		case currentSymbol == escape:
			if escaping {
				symbol = currentSymbol
				escaping = false
				continue
			}

			if symbol != none {
				sb.WriteRune(symbol)
				symbol = none
			}

			escaping = true

		case unicode.IsLetter(currentSymbol):
			if escaping {
				return "", ErrInvalidString
			}

			if symbol != none {
				sb.WriteRune(symbol)
			}

			symbol = currentSymbol
		}
	}

	if symbol != none {
		sb.WriteRune(symbol)
	}

	return sb.String(), nil
}
