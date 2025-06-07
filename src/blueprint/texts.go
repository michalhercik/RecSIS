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
	ttUnassignYear1    string
	ttUnassignYear2    string
	ttRemoveYear1      string
	ttRemoveYear2      string
	ttUnassignWinter   string
	ttRemoveWinter     string
	ttUnassign         string
	ttReassign         string
	ttUnassignSummer   string
	ttRemoveSummer     string
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
		ttUnassignYear1:    "Přesunout všechny předměty z ",
		ttUnassignYear2:    ". ročníku do nezařazených",
		ttRemoveYear1:      "Odstranit všechny předměty z ",
		ttRemoveYear2:      ". ročníku",
		ttUnassignWinter:   "Přesunout všechny předměty z tohoto zimního semestru do nezařazených",
		ttRemoveWinter:     "Odstranit všechny předměty z tohoto zimního semestru",
		ttUnassign:         "Přesunout předmět do nezařazených",
		ttReassign:         "Přesunout předmět",
		ttUnassignSummer:   "Přesunout všechny předměty z tohoto letního semestru do nezařazených",
		ttRemoveSummer:     "Odstranit všechny předměty z tohoto letního semestru",
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
		// language
		language: language.CS,
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
		ttUnassignYear1:    "Unassign all courses from Year ",
		ttUnassignYear2:    "",
		ttRemoveYear1:      "Remove all courses in Year ",
		ttRemoveYear2:      "",
		ttUnassignWinter:   "Unassign all courses from this winter semester",
		ttRemoveWinter:     "Remove all courses from this winter semester",
		ttUnassign:         "Unassign course",
		ttReassign:         "Reassign course",
		ttUnassignSummer:   "Unassign all courses from this summer semester",
		ttRemoveSummer:     "Remove all courses from this summer semester",
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
		// language
		language: language.EN,
	},
}
