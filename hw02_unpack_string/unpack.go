package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(v string) (string, error) {
	if v == "" {
		return "", nil
	}

	var baffle, number bool
	result := &strings.Builder{}
	smbRep := &strings.Builder{}

	if unicode.IsDigit(rune(v[0])) {
		return "", ErrInvalidString
	}

	for _, currRune := range v {
		if err := check(&baffle, &number, currRune, result, smbRep); err != nil {
			return "", err
		}
	}

	if baffle {
		return result.String(), ErrInvalidString
	}

	result.WriteString(smbRep.String())
	return result.String(), nil
}

func check(baffle, number *bool, currRune rune, result, smbRep *strings.Builder) error {
	if currRune == '\\' {
		*number = false
		if *baffle {
			if smbRep.Len() == 0 {
				smbRep.WriteRune('\\')
				*baffle = false
				return nil
			}
			result.WriteString(smbRep.String())
			smbRep.Reset()
			return nil
		}
		*baffle = true
		return nil
	}

	if num, err := strconv.Atoi(string(currRune)); err == nil {
		if *number {
			return ErrInvalidString
		}

		if *baffle {
			smbRep.WriteRune(currRune)
			*baffle = false
			return nil
		}

		*number = true
		if smbRep.Len() == 0 {
			buf := result.String()
			result.Reset()
			result.WriteString(buf[:len(buf)-1])
			smbRep.WriteRune(rune(buf[len(buf)-1]))
		}

		*baffle = false
		if num != 0 {
			result.WriteString(strings.Repeat(smbRep.String(), num))
		}
		smbRep.Reset()
		return nil
	}

	if *baffle {
		if smbRep.Len() == 0 {
			smbRep.WriteRune('\\')
			*baffle = false
		}

		smbRep.WriteRune(currRune)
		return nil
	}

	result.WriteRune(currRune)
	*number = false
	return nil
}
