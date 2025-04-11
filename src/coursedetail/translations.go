package coursedetail

import (
	"strconv"

	"github.com/michalhercik/RecSIS/language"
	"github.com/michalhercik/RecSIS/utils"
)

type text struct {
	Language         string
	Faculty          string
	Semester         string
	Winter           string
	Summer           string
	Both             string
	WinterAssign     string
	SummerAssign     string
	BothAssign       string
	Year             string
	Assign           string
	ECredits         string
	Capacity         string
	CapacityNoLimit  string
	ScopeExam        string
	StateOfCourse    string
	LanguageOfCourse string
	AdditionalInfo   string
	Prerequisites    string
	Corequisites     string
	Interchange      string
	Incompatible     string
	Classes          string
	Classification   string
	Guarantors       string
	Teachers         string
	ShowMore         string
	ShowLess         string
	NotRated         string
	NoRatings        string
	SISLink          string
	Detail           string
	Comments         string
	NoComments       string
	CourseWithCode   string
	NotFound         string
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
		Language:         "cs",
		Faculty:          "Fakulta",
		Semester:         "Semestr",
		Winter:           "Zimní",
		Summer:           "Letní",
		Both:             "Oba",
		WinterAssign:     "ZS",
		SummerAssign:     "LS",
		BothAssign:       "Oba",
		Year:             "ročník",
		Assign:           "Přiřadit",
		ECredits:         "E-Kredity",
		Capacity:         "Kapacita",
		CapacityNoLimit:  "bez omezení",
		ScopeExam:        "Rozsah, examinace",
		StateOfCourse:    "Stav předmětu",
		LanguageOfCourse: "Jazyk výuky",
		AdditionalInfo:   "Další informace",
		Prerequisites:    "Prerekvizity",
		Corequisites:     "Korekvizity",
		Interchange:      "Záměnnost",
		Incompatible:     "Neslučitelnost",
		Classes:          "Třída(y)",
		Classification:   "Klasifikace",
		Guarantors:       "Garant(i)",
		Teachers:         "Vyučující",
		ShowMore:         "Zobrazit více",
		ShowLess:         "Zobrazit méně",
		NotRated:         "Nehodnoceno",
		NoRatings:        "Žádná hodnocení",
		SISLink:          "Odkaz do SIS",
		Detail:           "Detailní informace",
		Comments:         "Komentáře",
		NoComments:       "Pro tento předmět nejsou žádné komentáře.",
		CourseWithCode:   "Předmět s kódem ",
		NotFound:         " nenalezen.",
		// utils
		Utils: utils.Texts["cs"],
	},
	language.EN: {
		Language:         "en",
		Faculty:          "Faculty",
		Semester:         "Semester",
		Winter:           "Winter",
		Summer:           "Summer",
		Both:             "Both",
		WinterAssign:     "Winter",
		SummerAssign:     "Summer",
		BothAssign:       "Both",
		Year:             "Year",
		Assign:           "Assign",
		ECredits:         "E-Credits",
		Capacity:         "Capacity",
		CapacityNoLimit:  "No limit",
		ScopeExam:        "Scope, examination",
		StateOfCourse:    "State of the course",
		LanguageOfCourse: "Language",
		AdditionalInfo:   "Additional information",
		Prerequisites:    "Pre-requisite",
		Corequisites:     "Co-requisite",
		Interchange:      "Interchangeability",
		Incompatible:     "Incompatibility",
		Classes:          "Class(es)",
		Classification:   "Classification",
		Guarantors:       "Guarantor(s)",
		Teachers:         "Teacher(s)",
		ShowMore:         "Show more",
		ShowLess:         "Show less",
		NotRated:         "Not rated",
		NoRatings:        "No ratings",
		SISLink:          "Link to SIS",
		Detail:           "Detailed information",
		Comments:         "Comments",
		NoComments:       "No comments for this course.",
		CourseWithCode:   "Course with code ",
		NotFound:         " not found.",
		// utils
		Utils: utils.Texts["en"],
	},
}
