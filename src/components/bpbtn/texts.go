package bpbtn

import (
	"strconv"

	"github.com/michalhercik/RecSIS/language"
)

type text struct {
	assign       string
	year         string
	winterAssign string
	summerAssign string
	// language
	language language.Language
}

func (t text) yearStr(year int) string {
	if t.language == language.CS {
		return strconv.Itoa(year) + ". " + t.year
	} else if t.language == language.EN {
		return t.year + " " + strconv.Itoa(year)
	}
	return ""
}

var texts = map[language.Language]text{
	language.CS: {
		assign:       "Přiřadit",
		year:         "ročník",
		winterAssign: "ZS",
		summerAssign: "LS",
		// language
		language: language.CS,
	},
	language.EN: {
		assign:       "Assign",
		year:         "Year",
		winterAssign: "Winter",
		summerAssign: "Summer",
		// language
		language: language.EN,
	},
}
