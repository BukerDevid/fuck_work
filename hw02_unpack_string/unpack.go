package hw02unpackstring

import (
	"errors"
	"strconv"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(v string) (string, error) {
	/*Check string and get size*/
	count, ok := CheckString(v)
	if !ok {
		return "", ErrInvalidString
	}

	resultBuffer := make([]rune, count)

	var symbol rune
	var limit int
	var baffle bool
	var idx int

	for _, rn := range v {

		if rn == rune('\\') {
			/*Check doubling baffle*/
			if baffle {
				baffle = false

				resultBuffer[idx] = rn
				idx++
				continue
			}

			baffle = true
			continue
		}

		if num, err := strconv.Atoi(string(rn)); err == nil {
			if baffle {
				/*If baffle is used. Number is one symbol*/
				baffle = false
				resultBuffer[idx] = rn
				idx++
				continue
			}

			/*If buffle not use and rune is first number in series (>5) ((num-1)+count)*/
			symbol = resultBuffer[idx-1] //get last symbol

			if num == 0 {
				/*Reset last simbol id number repeate equal zero*/
				idx--
				limit = 0

			} else {
				/*Limit is last symbol (one in result) + other*/
				limit = idx + (num - 1)
			}

			for ; idx < limit; idx++ {
				resultBuffer[idx] = symbol
			}
			continue
		}

		resultBuffer[idx] = rn
		idx++
	}

	return string(resultBuffer), nil
}

func CheckString(v string) (int, bool) {
	if v == "" {
		return 0, true
	}

	/*First rune is number*/
	if _, err := strconv.Atoi(string([]rune(v)[0])); err == nil {
		return 0, false
	}

	var baffle bool = false
	var number bool = false
	var count int = 0

	for _, rn := range v {

		if rn == rune('\\') {
			/*Check doubling baffle*/
			if baffle {
				baffle = false
				number = false
				count++
				continue
			}

			baffle = true
			continue
		}

		if num, err := strconv.Atoi(string(rn)); err == nil {
			if baffle {
				/*If baffle is used. Number is one symbol*/
				baffle = false
				count++
				continue
			}

			/*If buffle not use but two number in series (5 >5) FALSE*/
			if number {
				return count, false
			}

			/*If buffle not use and rune is first number in series (>5) ((num-1)+count)*/
			count = (num - 1) + count
			number = true
			continue
		}

		if baffle {
			return count, false
		}

		number = false
		count++
	}

	if baffle {
		return count, false
	}

	return count, true
}
