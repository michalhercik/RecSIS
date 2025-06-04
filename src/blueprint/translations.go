package blueprint

import (
	"strconv"

	"github.com/michalhercik/RecSIS/language"
	"github.com/michalhercik/RecSIS/utils"
)

type text struct {
	PageTitle        string
	Total            string
	Code             string
	Title            string
	Credits          string
	CreditsShort     string
	Winter           string
	WinterLong       string
	Summer           string
	SummerLong       string
	Both             string
	Guarantors       string
	Unassigned       string
	NoUnassignedText string
	EmptySemester    string
	Year             string
	YearBig          string
	Semester         string
	NumOfYears       string
	// modal
	ModalTitle      string
	ModalContent    string
	Cancel          string
	RemoveCourses   string
	UnassignCourses string
	// tooltips
	TTUncheckAll       string
	TTUnassignChecked  string
	TTAssignChecked    string
	TTRemoveChecked    string
	TTRemoveUnassigned string
	TTAssign           string
	TTRemove           string
	TTMove             string
	TTUnassignYear1    string
	TTUnassignYear2    string
	TTRemoveYear1      string
	TTRemoveYear2      string
	TTUnassignWinter   string
	TTRemoveWinter     string
	TTUnassign         string
	TTReassign         string
	TTUnassignSummer   string
	TTRemoveSummer     string
	// stats
	TTNumberOfCredits string
	TTSemesterCredits string
	TTYearCredits     string
	TTRunningCredits  string
	TTTotalCredits    string
	// warnings
	WWrongAssignWinter    string
	WWrongAssignSummer    string
	WAssignedMoreThanOnce string
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
		PageTitle:        "Blueprint",
		Total:            "Celkem",
		Code:             "Kód",
		Title:            "Název",
		Credits:          "Kredity",
		CreditsShort:     "Kr.",
		Winter:           "ZS",
		WinterLong:       "Zimní",
		Summer:           "LS",
		SummerLong:       "Letní",
		Both:             "Oba",
		Guarantors:       "Garant(i)",
		Unassigned:       "Nezařazené",
		NoUnassignedText: "Žádné nezařazené předměty",
		EmptySemester:    "žádné předměty",
		Year:             "ročník",
		YearBig:          "Ročník",
		Semester:         "Semestr",
		NumOfYears:       "Počet ročníků",
		// modal
		ModalTitle:      "Chystáte se odstranit neprázdný ročník",
		ModalContent:    "Ročník, který se chystáte odstranit, obsahuje předměty. Můžete je odstranit, přesunout do nezařazených nebo tuto akci zrušit.",
		Cancel:          "Zrušit",
		RemoveCourses:   "Odstranit",
		UnassignCourses: "Přesunout do nezařazených",
		// tooltips
		TTUncheckAll:       "Odznačit všechny označené předměty",
		TTUnassignChecked:  "Přesunout vybrané předměty do nezařazených",
		TTAssignChecked:    "Zařadit vybrané předměty",
		TTRemoveChecked:    "Odstranit vybrané předměty",
		TTRemoveUnassigned: "Odstranit nezařazené předměty",
		TTAssign:           "Zařadit předmět",
		TTRemove:           "Odstranit předmět",
		TTMove:             "Přesunout předmět pomocí drag-and-drop",
		TTUnassignYear1:    "Přesunout všechny předměty z ",
		TTUnassignYear2:    ". ročníku do nezařazených",
		TTRemoveYear1:      "Odstranit všechny předměty z ",
		TTRemoveYear2:      ". ročníku",
		TTUnassignWinter:   "Přesunout všechny předměty z tohoto zimního semestru do nezařazených",
		TTRemoveWinter:     "Odstranit všechny předměty z tohoto zimního semestru",
		TTUnassign:         "Přesunout předmět do nezařazených",
		TTReassign:         "Přesunout předmět",
		TTUnassignSummer:   "Přesunout všechny předměty z tohoto letního semestru do nezařazených",
		TTRemoveSummer:     "Odstranit všechny předměty z tohoto letního semestru",
		// stats
		TTNumberOfCredits: "Počet kreditů",
		TTSemesterCredits: "V tomto semestru",
		TTYearCredits:     "V tomto ročníku",
		TTRunningCredits:  "Průběžný součet",
		TTTotalCredits:    "Celkem",
		// warnings
		WWrongAssignWinter:    "Předmět je zařazen do zimního semestru, ale měl by být v letním semestru.",
		WWrongAssignSummer:    "Předmět je zařazen do letního semestru, ale měl by být v zimním semestru.",
		WAssignedMoreThanOnce: "Předmět je zařazen více než jednou ",
		// utils
		Utils: utils.Texts[language.CS],
	},
	language.EN: {
		PageTitle:        "Blueprint",
		Total:            "Total",
		Code:             "Code",
		Title:            "Title",
		Credits:          "Credits",
		CreditsShort:     "Cr.",
		Winter:           "Winter",
		WinterLong:       "Winter",
		Summer:           "Summer",
		SummerLong:       "Summer",
		Both:             "Both",
		Guarantors:       "Guarantor(s)",
		Unassigned:       "Unassigned",
		NoUnassignedText: "No unassigned courses",
		EmptySemester:    "no courses",
		Year:             "Year",
		YearBig:          "Year",
		Semester:         "Semester",
		NumOfYears:       "Number of years",
		// modal
		ModalTitle:      "You are about to remove a non-empty year",
		ModalContent:    "The year you are about to remove contains courses. You can remove them, unassign them or cancel this action.",
		Cancel:          "Cancel",
		RemoveCourses:   "Remove",
		UnassignCourses: "Unassign",
		// tooltips
		TTUncheckAll:       "Uncheck all selected courses",
		TTUnassignChecked:  "Unassign all selected courses",
		TTAssignChecked:    "Assign all selected courses",
		TTRemoveChecked:    "Remove all selected courses",
		TTRemoveUnassigned: "Remove all unassigned courses",
		TTAssign:           "Assign course",
		TTRemove:           "Remove course",
		TTMove:             "Drag and drop to sort",
		TTUnassignYear1:    "Unassign all courses from Year ",
		TTUnassignYear2:    "",
		TTRemoveYear1:      "Remove all courses in Year ",
		TTRemoveYear2:      "",
		TTUnassignWinter:   "Unassign all courses from this winter semester",
		TTRemoveWinter:     "Remove all courses from this winter semester",
		TTUnassign:         "Unassign course",
		TTReassign:         "Reassign course",
		TTUnassignSummer:   "Unassign all courses from this summer semester",
		TTRemoveSummer:     "Remove all courses from this summer semester",
		// stats
		TTNumberOfCredits: "Number of credits",
		TTSemesterCredits: "In this semester",
		TTYearCredits:     "In this year",
		TTRunningCredits:  "Running total",
		TTTotalCredits:    "Total",
		// warnings
		WWrongAssignWinter:    "Course is assigned in a winter semester (should be in summer).",
		WWrongAssignSummer:    "Course is assigned in a summer semester (should be in winter).",
		WAssignedMoreThanOnce: "Course is assigned more than once ",
		// utils
		Utils: utils.Texts[language.EN],
	},
}
