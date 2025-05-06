package bpbtn

import (
	"strconv"

	"github.com/michalhercik/RecSIS/language"
	"github.com/michalhercik/RecSIS/utils"
)

type Text struct {
	Language     string
	Assign       string
	Year         string
	WinterAssign string
	SummerAssign string
	Utils        utils.Text
}

func (t Text) YearStr(year int) string {
	if t.Language == "cs" {
		return strconv.Itoa(year) + ". " + t.Year
	} else if t.Language == "en" {
		return t.Year + " " + strconv.Itoa(year)
	}
	return ""
}

var texts = map[language.Language]Text{
	language.CS: {
		Language:     "cs",
		Assign:       "Přiřadit",
		Year:         "ročník",
		WinterAssign: "ZS",
		SummerAssign: "LS",
		Utils:        utils.Texts["cs"],
	},
	language.EN: {
		Language:     "en",
		Assign:       "Assign",
		Year:         "Year",
		WinterAssign: "Winter",
		SummerAssign: "Summer",
		Utils:        utils.Texts["en"],
	},
}
