package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(v string) (string, error) {
	result := &strings.Builder{}
	if v == "" {
		return "", nil
	}

	if unicode.IsDigit(rune(v[0])) {
		return "", ErrInvalidString
	}

	sym_rep := &strings.Builder{}
	var baffle bool

	for _, cur_run := range v {
		if cur_run == '\\' {
			if baffle {
				baffle = false
				sym_rep.WriteByte(byte(cur_run))
				// must_num = true
				continue
			}

			sym_rep = &strings.Builder{}
		}

		if num, err := strconv.Atoi(string(cur_run)); err == nil {
			if baffle && sym_rep.Len() == 0 {
				sym_rep.WriteByte(byte(cur_run))
				// must_num = true
				continue
			}

			if baffle {
				baffle = false
				result.WriteString(strings.Repeat(sym_rep.String(), num-1))
				continue
			}

			result.WriteString(strings.Repeat(result.String()[result.Len():], num-1))
			continue
		}

		if baffle {
			sym_rep.WriteByte(byte(cur_run))
			continue
		}

		result.WriteByte(byte(cur_run))
	}

	if baffle {
		return "", ErrInvalidString
	}

	return result.String(), nil
}
