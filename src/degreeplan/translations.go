package degreeplan

import (
	"strconv"

	"github.com/michalhercik/RecSIS/language"
	"github.com/michalhercik/RecSIS/utils"
)

type text struct {
	Language     string
	Code         string
	Title        string
	Status       string
	Completed    string
	InBlueprint  string
	NotCompleted string
	Credits      string
	CreditsShort string
	Needed       string
	Winter       string
	Summer       string
	Blueprint    string
	Assign       string
	Year         string
	WinterAssign string
	SummerAssign string
	// utils
	Utils utils.Text
}

func (t text) YearStr(year int) string {
	if t.Language == "cs" {
		return strconv.Itoa(year) + ". " + t.Year
	} else if t.Language == "en" {
		return t.Year + " " + strconv.Itoa(year)
	}
	return ""
}

var texts = map[language.Language]text{
	language.CS: {
		Language:     "cs",
		Code:         "Kód",
		Title:        "Název",
		Status:       "Stav",
		Completed:    "Splněno",
		InBlueprint:  "Blueprint",
		NotCompleted: "Nesplněno",
		Credits:      "Kredity",
		CreditsShort: "kr.",
		Needed:       "potřeba",
		Winter:       "ZS",
		Summer:       "LS",
		Blueprint:    "Blueprint",
		Assign:       "Přiřadit",
		Year:         "ročník",
		WinterAssign: "ZS",
		SummerAssign: "LS",
		// utils
		Utils: utils.Texts["cs"],
	},
	language.EN: {
		Language:     "en",
		Code:         "Code",
		Title:        "Title",
		Status:       "Status",
		Completed:    "Completed",
		InBlueprint:  "Blueprint",
		NotCompleted: "Not completed",
		Credits:      "Credits",
		CreditsShort: "cr.",
		Needed:       "needed",
		Winter:       "Winter",
		Summer:       "Summer",
		Blueprint:    "Blueprint",
		Assign:       "Assign",
		Year:         "Year",
		WinterAssign: "Winter",
		SummerAssign: "Summer",
		// utils
		Utils: utils.Texts["en"],
	},
}
