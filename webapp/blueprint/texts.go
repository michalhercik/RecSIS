package blueprint

import (
	"strconv"

	"github.com/michalhercik/RecSIS/language"
)

type text struct {
	pageTitle        string
	total            string
	code             string
	title            string
	credits          string
	creditsShort     string
	winter           string
	winterLong       string
	summer           string
	summerLong       string
	both             string
	guarantors       string
	unassigned       string
	noUnassignedText string
	emptySemester    string
	year             string
	yearBig          string
	semester         string
	numOfYears       string
	// modal
	modalTitle      string
	modalContent    string
	cancel          string
	removeCourses   string
	unassignCourses string
	// tooltips
	ttUncheckAll       string
	ttUnassignChecked  string
	ttAssignChecked    string
	ttRemoveChecked    string
	ttRemoveUnassigned string
	ttAssign           string
	ttRemove           string
	ttMove             string
	ttUnassignWinter   string
	ttRemoveWinter     string
	ttUnassign         string
	ttReassign         string
	ttUnassignSummer   string
	ttRemoveSummer     string
	ttAssignedCredits  string
	ttBlueprintCredits string
	// stats
	ttNumberOfCredits string
	ttSemesterCredits string
	ttYearCredits     string
	ttRunningCredits  string
	ttTotalCredits    string
	// warnings
	wWrongAssignWinter    string
	wWrongAssignSummer    string
	wAssignedMoreThanOnce string
	wPrerequisiteNotMet   string
	wCorequisiteNotMet    string
	wIncompatiblePresent  string
	// language
	language language.Language
	// errors
	errCannotGetBlueprint             string
	errNoSemestersFound               string
	errInvalidYearInDB                string
	errInvalidSemesterInDB            string
	errDuplicateCourseInBP            string
	errDuplicateCoursesInBP           string
	errDuplicateCoursesInBPUnassigned string
	errCannotMoveCourses              string
	errCannotAppendCourses            string
	errCannotUnassignYear             string
	errCannotUnassignSemester         string
	errCannotRemoveCourses            string
	errCannotAddYear                  string
	errCannotRemoveYear               string
	errCannotUnFoldSemester           string
	errMissingCourseID                string
	errInvalidCourseID                string
	errInvalidMoveType                string
	errInvalidRemoveType              string
	errMissingUnassignParam           string
	errInvalidUnassignParam           string
	errMissingFoldedParam             string
	errInvalidFoldedParam             string
	errMissingYearParam               string
	errInvalidYearParam               string
	errMissingSemesterParam           string
	errInvalidSemesterParam           string
	errMissingPositionParam           string
	errInvalidPositionParam           string
	errPageNotFound                   string
}

func (t text) yearStr(year int) string {
	switch t.language {
	case language.CS:
		return strconv.Itoa(year) + ". " + t.year
	case language.EN:
		return t.year + " " + strconv.Itoa(year)
	}
	return ""
}

var texts = map[language.Language]text{
	language.CS: {
		pageTitle:        "Blueprint",
		total:            "Celkem",
		code:             "Kód",
		title:            "Název",
		credits:          "Kredity",
		creditsShort:     "Kr.",
		winter:           "ZS",
		winterLong:       "Zimní",
		summer:           "LS",
		summerLong:       "Letní",
		both:             "Oba",
		guarantors:       "Garant(i)",
		unassigned:       "Nezařazené",
		noUnassignedText: "Žádné nezařazené předměty",
		emptySemester:    "žádné předměty",
		year:             "ročník",
		yearBig:          "Ročník",
		semester:         "Semestr",
		numOfYears:       "Počet ročníků",
		// modal
		modalTitle:      "Chystáte se odstranit neprázdný ročník",
		modalContent:    "Ročník, který se chystáte odstranit, obsahuje předměty. Můžete je odstranit, přesunout do nezařazených nebo tuto akci zrušit.",
		cancel:          "Zrušit",
		removeCourses:   "Odstranit",
		unassignCourses: "Přesunout do nezařazených",
		// tooltips
		ttUncheckAll:       "Odznačit všechny označené předměty",
		ttUnassignChecked:  "Přesunout vybrané předměty do nezařazených",
		ttAssignChecked:    "Zařadit vybrané předměty",
		ttRemoveChecked:    "Odstranit vybrané předměty",
		ttRemoveUnassigned: "Odstranit nezařazené předměty",
		ttAssign:           "Zařadit předmět",
		ttRemove:           "Odstranit předmět",
		ttMove:             "Přesunout předmět pomocí drag-and-drop",
		ttUnassignWinter:   "Přesunout všechny předměty z tohoto zimního semestru do nezařazených",
		ttRemoveWinter:     "Odstranit všechny předměty z tohoto zimního semestru",
		ttUnassign:         "Přesunout předmět do nezařazených",
		ttReassign:         "Přesunout předmět",
		ttUnassignSummer:   "Přesunout všechny předměty z tohoto letního semestru do nezařazených",
		ttRemoveSummer:     "Odstranit všechny předměty z tohoto letního semestru",
		ttAssignedCredits:  "Počet kreditů přiřazených do ročníků",
		ttBlueprintCredits: "Počet kreditů v blueprintu",
		// stats
		ttNumberOfCredits: "Počet kreditů",
		ttSemesterCredits: "V tomto semestru",
		ttYearCredits:     "V tomto ročníku",
		ttRunningCredits:  "Průběžný součet",
		ttTotalCredits:    "Celkem",
		// warnings
		wWrongAssignWinter:    "Předmět je zařazen do zimního semestru, ale měl by být v letním semestru.",
		wWrongAssignSummer:    "Předmět je zařazen do letního semestru, ale měl by být v zimním semestru.",
		wAssignedMoreThanOnce: "Předmět je zařazen více než jednou ",
		wPrerequisiteNotMet:   "Nesplněna prerekvizita pro tento předmět.",
		wCorequisiteNotMet:    "Nesplněna korekvizita pro tento předmět.",
		wIncompatiblePresent:  "V blueprintu je přítomen neslučitelný předmět.",
		// language
		language: language.CS,
		// errors
		errCannotGetBlueprint:             "Nelze načíst stránku Blueprint z databáze",
		errNoSemestersFound:               "Pro uživatele nebyl nalezen žádný semestr v databázi",
		errInvalidYearInDB:                "Neplatný ročník v databázi",
		errInvalidSemesterInDB:            "Neplatný semestr v databázi",
		errDuplicateCourseInBP:            "Vybraný kurz je již ve vybraném ročníku a semestru v Blueprintu přiřazen",
		errDuplicateCoursesInBP:           "Jeden nebo více zvolených kurzů jsou již ve vybraném ročníku a semestru v Blueprintu přiřazeny",
		errDuplicateCoursesInBPUnassigned: "Jeden nebo více zvolených kurzů jsou již v nezařazených kurzech v Blueprintu",
		errCannotMoveCourses:              "Nelze přesunout vybrané kurzy",
		errCannotAppendCourses:            "Nelze přesunout vybrané kurzy",
		errCannotUnassignYear:             "Nelze přesunout kurzy z ročníku do nezařazených",
		errCannotUnassignSemester:         "Nelze přesunout kurzy z tohoto semestru do nezařazených",
		errCannotRemoveCourses:            "Nelze odstranit vybrané kurzy",
		errCannotAddYear:                  "Nelze přidat ročník",
		errCannotRemoveYear:               "Nelze odstranit ročník",
		errCannotUnFoldSemester:           "Nelze sbalit/rozbalit semestr",
		errMissingCourseID:                "Chybí ID kurzu",
		errInvalidCourseID:                "Neplatné ID kurzu (musí být celé kladné číslo)",
		errInvalidMoveType:                "Neplatný typ přesunu předmětu",
		errInvalidRemoveType:              "Neplatný typ odstranění předmětu",
		errMissingUnassignParam:           "Chybí parametr jestli přesunout předměty do nezařazených",
		errInvalidUnassignParam:           "Neplatný parametr přesunu předmětů do nezařazených (musí být true nebo false)",
		errMissingFoldedParam:             "Chybí parametr jestli je semestr sbalený nebo rozbalený",
		errInvalidFoldedParam:             "Neplatný parametr sbalení/rozbalení semestru (musí být true nebo false)",
		errMissingYearParam:               "Chybí parametr ročníku",
		errInvalidYearParam:               "Neplatný parametr ročníku (musí být celé kladné číslo větší nebo rovno 0)",
		errMissingSemesterParam:           "Chybí parametr semestru",
		errInvalidSemesterParam:           "Neplatný parametr semestru",
		errMissingPositionParam:           "Chybí parametr pozice",
		errInvalidPositionParam:           "Neplatný parametr pozice (musí být celé kladné číslo větší nebo rovno -1)",
		errPageNotFound:                   "Stránka nenalezena",
	},
	language.EN: {
		pageTitle:        "Blueprint",
		total:            "Total",
		code:             "Code",
		title:            "Title",
		credits:          "Credits",
		creditsShort:     "Cr.",
		winter:           "Winter",
		winterLong:       "Winter",
		summer:           "Summer",
		summerLong:       "Summer",
		both:             "Both",
		guarantors:       "Guarantor(s)",
		unassigned:       "Unassigned",
		noUnassignedText: "No unassigned courses",
		emptySemester:    "no courses",
		year:             "Year",
		yearBig:          "Year",
		semester:         "Semester",
		numOfYears:       "Number of years",
		// modal
		modalTitle:      "You are about to remove a non-empty year",
		modalContent:    "The year you are about to remove contains courses. You can remove them, unassign them or cancel this action.",
		cancel:          "Cancel",
		removeCourses:   "Remove",
		unassignCourses: "Unassign",
		// tooltips
		ttUncheckAll:       "Uncheck all selected courses",
		ttUnassignChecked:  "Unassign all selected courses",
		ttAssignChecked:    "Assign all selected courses",
		ttRemoveChecked:    "Remove all selected courses",
		ttRemoveUnassigned: "Remove all unassigned courses",
		ttAssign:           "Assign course",
		ttRemove:           "Remove course",
		ttMove:             "Drag and drop to sort",
		ttUnassignWinter:   "Unassign all courses from this winter semester",
		ttRemoveWinter:     "Remove all courses from this winter semester",
		ttUnassign:         "Unassign course",
		ttReassign:         "Reassign course",
		ttUnassignSummer:   "Unassign all courses from this summer semester",
		ttRemoveSummer:     "Remove all courses from this summer semester",
		ttAssignedCredits:  "Number of credits assigned to years",
		ttBlueprintCredits: "Number of credits in the blueprint",
		// stats
		ttNumberOfCredits: "Number of credits",
		ttSemesterCredits: "In this semester",
		ttYearCredits:     "In this year",
		ttRunningCredits:  "Running total",
		ttTotalCredits:    "Total",
		// warnings
		wWrongAssignWinter:    "Course is assigned in a winter semester (should be in summer).",
		wWrongAssignSummer:    "Course is assigned in a summer semester (should be in winter).",
		wAssignedMoreThanOnce: "Course is assigned more than once ",
		wPrerequisiteNotMet:   "Prerequisite for this course is not met.",
		wCorequisiteNotMet:    "Corequisite for this course is not met.",
		wIncompatiblePresent:  "An incompatible course is present in the blueprint.",
		// language
		language: language.EN,
		// errors
		errCannotGetBlueprint:             "Unable to retrieve Blueprint page from database",
		errNoSemestersFound:               "No semesters found for user in database",
		errInvalidYearInDB:                "Invalid year in database",
		errInvalidSemesterInDB:            "Invalid semester in database",
		errDuplicateCourseInBP:            "Selected course is already assigned in the selected year and semester in Blueprint",
		errDuplicateCoursesInBP:           "One or more selected courses are already assigned in the selected year and semester in Blueprint",
		errDuplicateCoursesInBPUnassigned: "One or more selected courses are already in the unassigned courses in Blueprint",
		errCannotMoveCourses:              "Cannot move selected courses",
		errCannotAppendCourses:            "Cannot append selected courses",
		errCannotUnassignYear:             "Cannot unassign courses from this year",
		errCannotUnassignSemester:         "Cannot unassign courses from this semester",
		errCannotRemoveCourses:            "Cannot remove selected courses",
		errCannotAddYear:                  "Cannot add year",
		errCannotRemoveYear:               "Cannot remove year",
		errCannotUnFoldSemester:           "Cannot fold/unfold semester",
		errMissingCourseID:                "Missing course ID",
		errInvalidCourseID:                "Invalid course ID (must be a positive integer)",
		errInvalidMoveType:                "Invalid course move type",
		errInvalidRemoveType:              "Invalid course remove type",
		errMissingUnassignParam:           "Missing parameter if should unassign courses",
		errInvalidUnassignParam:           "Invalid unassign parameter (must be true or false)",
		errMissingFoldedParam:             "Missing parameter if semester is folded or unfolded",
		errInvalidFoldedParam:             "Invalid folded parameter (must be true or false)",
		errMissingYearParam:               "Missing year parameter",
		errInvalidYearParam:               "Invalid year parameter (must be a positive integer greater than or equal to 0)",
		errMissingSemesterParam:           "Missing semester parameter",
		errInvalidSemesterParam:           "Invalid semester parameter",
		errMissingPositionParam:           "Missing position parameter",
		errInvalidPositionParam:           "Invalid position parameter (must be a positive integer greater than or equal to -1)",
		errPageNotFound:                   "Page not found",
	},
}
