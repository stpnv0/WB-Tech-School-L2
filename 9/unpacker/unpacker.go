package unpacker

import (
	"errors"
	"strconv"
	"unicode"
)

var (
	firstRuneIsDigitErr = errors.New("first character is a digit")
	trailingEscapeErr   = errors.New("string ends with escape character")
	invalidSequenceErr  = errors.New("invalid sequence")
	multiDigitNumberErr = errors.New("multi-digit numbers are not supported")
)

func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	runes := []rune(s)
	if unicode.IsNumber(runes[0]) {
		return "", firstRuneIsDigitErr
	}

	escape := false
	var result []rune
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if escape {
			result = append(result, r)
			escape = false
			continue
		}

		if r == '\\' {
			if i+1 >= len(runes) {
				return "", trailingEscapeErr
			}

			escape = true
			continue
		}

		if unicode.IsDigit(r) {
			if len(result) == 0 {
				return "", invalidSequenceErr
			}
			if i+1 < len(runes) && unicode.IsDigit(runes[i+1]) {
				return "", multiDigitNumberErr
			}

			count, err := strconv.Atoi(string(r))
			if err != nil {
				return "", err
			}

			symbolRepeat := result[len(result)-1]
			if unicode.IsSpace(symbolRepeat) {
				return "", invalidSequenceErr
			}
			result = result[:len(result)-1]

			if count > 0 {
				for k := 0; k < count; k++ {
					result = append(result, symbolRepeat)
				}
			}
			continue
		}

		result = append(result, r)
	}

	if escape {
		return "", trailingEscapeErr
	}

	return string(result), nil
}
