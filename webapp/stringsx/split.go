package stringsx

import (
	"strings"
)

func SplitByLastSpace(s string) (string, string) {
	words := strings.Fields(s)
	rest, last := "", ""
	if len(words) > 1 {
		rest = strings.Join(words[:len(words)-1], " ")
		last = " " + words[len(words)-1]
	} else {
		rest = ""
		last = s
	}
	return rest, last
}
