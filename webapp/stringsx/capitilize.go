package stringsx

import (
	"strings"
	"unicode/utf8"
)

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	firstRune, _ := utf8.DecodeRuneInString(s)
	return strings.ToUpper(string(firstRune)) + s[1:]
}
