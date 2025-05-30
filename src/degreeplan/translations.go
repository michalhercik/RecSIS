package degreeplan

import (
	"strconv"

	"github.com/michalhercik/RecSIS/language"
	"github.com/michalhercik/RecSIS/utils"
)

type text struct {
	PageTitle          string
	Code               string
	Title              string
	Status             string
	Completed          string
	InBlueprint        string
	Unassigned         string
	NotCompleted       string
	Credits            string
	CreditsShort       string
	Needed             string
	Winter             string
	Summer             string
	Both               string
	Guarantors         string
	Blueprint          string
	Assign             string
	CourseIsUnassigned string
	Year               string
	WinterAssign       string
	SummerAssign       string
	// utils
	Utils utils.Text
}

func (t text) YearStr(year int) string {
	if t.Utils.Language == language.CS {
		return strconv.Itoa(year) + ". " + t.Year
	} else if t.Utils.Language == language.EN {
		return t.Year + " " + strconv.Itoa(year)
	}
	return ""
}

var texts = map[language.Language]text{
	language.CS: {
		PageTitle:          "Studijní plán",
		Code:               "Kód",
		Title:              "Název",
		Status:             "Stav",
		Completed:          "Splněno",
		InBlueprint:        "Blueprint",
		Unassigned:         "Nezařazen",
		NotCompleted:       "Nesplněno",
		Credits:            "Kredity",
		CreditsShort:       "Kr.",
		Needed:             "potřeba",
		Winter:             "ZS",
		Summer:             "LS",
		Both:               "Oba",
		Guarantors:         "Garant(i)",
		Blueprint:          "Blueprint",
		Assign:             "Přiřadit",
		CourseIsUnassigned: "Kurz je v blueprintu, ale není zařazen.",
		Year:               "ročník",
		WinterAssign:       "ZS",
		SummerAssign:       "LS",
		// utils
		Utils: utils.Texts[language.CS],
	},
	language.EN: {
		PageTitle:          "Degree Plan",
		Code:               "Code",
		Title:              "Title",
		Status:             "Status",
		Completed:          "Completed",
		InBlueprint:        "Blueprint",
		Unassigned:         "Unassigned",
		NotCompleted:       "Not completed",
		Credits:            "Credits",
		CreditsShort:       "Cr.",
		Needed:             "needed",
		Winter:             "Winter",
		Both:               "Both",
		Summer:             "Summer",
		Guarantors:         "Guarantor(s)",
		Blueprint:          "Blueprint",
		Assign:             "Assign",
		CourseIsUnassigned: "Course is in the blueprint but not assigned.",
		Year:               "Year",
		WinterAssign:       "Winter",
		SummerAssign:       "Summer",
		// utils
		Utils: utils.Texts[language.EN],
	},
}
