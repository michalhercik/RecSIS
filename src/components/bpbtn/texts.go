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
	// errors
	errNoCoursesProvided    string
	errNoYearProvided       string
	errInvalidYear          string
	errNoSemesterProvided   string
	errInvalidSemester      string
	errDuplicateCourseInBP  string
	errDuplicateCoursesInBP string
	errAddCourseToBPFailed  string
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
		// errors
		errNoCoursesProvided:    "Nebyl poskytnut žádný kurz pro přiřazení do Blueprintu",
		errNoYearProvided:       "Nebyl poskytnut ročník pro přiřazení do Blueprintu",
		errInvalidYear:          "Poskytnutý ročník není platný pro přiřazení do Blueprintu",
		errNoSemesterProvided:   "Nebyl poskytnut semestr pro přiřazení do Blueprintu",
		errInvalidSemester:      "Poskytnutý semestr není platný pro přiřazení do Blueprintu",
		errDuplicateCourseInBP:  "Vybraný kurz je již ve vybraném ročníku a semestru v Blueprintu přiřazen",
		errDuplicateCoursesInBP: "Jeden nebo více zvolených kurzů jsou již ve vybraném ročníku a semestru v Blueprintu přiřazeny",
		errAddCourseToBPFailed:  "Přidání kurzu/ů do Blueprintu selhalo",
	},
	language.EN: {
		assign:       "Assign",
		year:         "Year",
		winterAssign: "Winter",
		summerAssign: "Summer",
		// language
		language: language.EN,
		// errors
		errNoCoursesProvided:    "No courses provided for Blueprint assignment",
		errNoYearProvided:       "No year provided for Blueprint assignment",
		errInvalidYear:          "Provided year is not valid for Blueprint assignment",
		errNoSemesterProvided:   "No semester provided for Blueprint assignment",
		errInvalidSemester:      "Provided semester is not valid for Blueprint assignment",
		errDuplicateCourseInBP:  "Selected course is already assigned in the selected year and semester in Blueprint",
		errDuplicateCoursesInBP: "One or more selected courses are already assigned in the selected year and semester in Blueprint",
		errAddCourseToBPFailed:  "Failed to add course(s) to Blueprint",
	},
}
