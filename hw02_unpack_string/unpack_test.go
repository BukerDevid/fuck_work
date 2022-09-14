package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

func TestCheckString(t *testing.T) {
	strings := []struct {
		input  string
		count  int
		result bool
	}{
		{input: "a4bc2d5e", count: 13, result: true},
		{input: "abccd", count: 5, result: true},
		{input: "", count: 0, result: true},
		{input: "aaa0b", count: 3, result: true},
		{input: "3abc", count: 0, result: false},
		{input: "45", count: 0, result: false},
		{input: "aaa10b", count: 3, result: false},
		{input: `qwe\4\5`, count: 5, result: true},
		{input: `qwe\45`, count: 8, result: true},
		{input: `qwe\\5`, count: 8, result: true},
		{input: `qwe\\\3`, count: 5, result: true},
		{input: `qwe\a`, count: 0, result: false},
		{input: `qwe\-5a`, count: 0, result: false},
	}

	for _, tc := range strings {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			count, ok := CheckString(tc.input)
			require.Truef(t, tc.result == ok, "error ok. Want ok %v, have ok %v", tc.result, ok)
			if tc.result {
				require.Truef(t, tc.count == count, "error count. Want count %d, have count %d", tc.count, count)
			}
		})
	}
}
