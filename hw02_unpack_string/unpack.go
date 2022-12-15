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
		if currRune == '\\' {
			number = false
			if baffle {
				if smbRep.Len() == 0 {
					smbRep.WriteRune('\\')
					baffle = false
					continue
				}
				result.WriteString(smbRep.String())
				smbRep.Reset()
				continue
			}
			baffle = true
			continue
		}
		if num, err := strconv.Atoi(string(currRune)); err == nil {
			if number {
				return result.String(), ErrInvalidString
			}
			if baffle {
				smbRep.WriteRune(currRune)
				baffle = false
				continue
			}
			number = true
			if smbRep.Len() == 0 {
				buf := result.String()
				result.Reset()
				result.WriteString(buf[:len(buf)-1])
				smbRep.WriteRune(rune(buf[len(buf)-1]))
			}
			baffle = false
			if num == 0 {
				smbRep.Reset()
				continue
			}
			result.WriteString(strings.Repeat(smbRep.String(), num))
			smbRep.Reset()
			continue
		}
		if baffle {
			if smbRep.Len() == 0 {
				smbRep.WriteRune('\\')
				baffle = false
			}
			smbRep.WriteRune(currRune)
		} else {
			result.WriteRune(currRune)
		}
		number = false
	}
	if baffle {
		return result.String(), ErrInvalidString
	}
	result.WriteString(smbRep.String())
	return result.String(), nil
}
