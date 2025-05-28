package coursedetail

import (
	"strconv"

	"github.com/michalhercik/RecSIS/language"
	"github.com/michalhercik/RecSIS/utils"
)

type text struct {
	Faculty                 string
	Department              string
	Semester                string
	Winter                  string
	Summer                  string
	Both                    string
	WinterAssign            string
	SummerAssign            string
	BothAssign              string
	Unassigned              string
	Year                    string
	Assign                  string
	ECredits                string
	InDegreePlan            string
	Completed               string
	Capacity                string
	CapacityNoLimit         string
	ScopeExam               string
	StateOfCourse           string
	LanguageOfCourse        string
	AdditionalInfo          string
	Prerequisites           string
	Corequisites            string
	Interchange             string
	Incompatible            string
	Classes                 string
	Classification          string
	Guarantors              string
	Teachers                string
	CourseGlobalNotes       string
	NotRated                string
	NoRatings               string
	AdditionalRatings       string
	CategoricalRatings      string
	SISLink                 string
	Description             string
	NoDescription           string
	Survey                  string
	NoSurvey                string
	SurveySearchPlaceholder string
	ScrollToSearch          string
	Detail                  string
	NoDetail                string
	CourseWithCode          string
	NotFound                string
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
		Faculty:                 "Fakulta",
		Department:              "Katedra",
		Semester:                "Semestr",
		Winter:                  "Zimní",
		Summer:                  "Letní",
		Both:                    "Oba",
		WinterAssign:            "ZS",
		SummerAssign:            "LS",
		BothAssign:              "Oba",
		Unassigned:              "Nezařazen",
		Year:                    "ročník",
		Assign:                  "Přiřadit",
		ECredits:                "E-Kredity",
		InDegreePlan:            "V studijním plánu",
		Completed:               "Splněno",
		Capacity:                "Kapacita",
		CapacityNoLimit:         "bez omezení",
		ScopeExam:               "Rozsah, examinace",
		StateOfCourse:           "Stav předmětu",
		LanguageOfCourse:        "Jazyk výuky",
		AdditionalInfo:          "Další informace",
		Prerequisites:           "Prerekvizity",
		Corequisites:            "Korekvizity",
		Interchange:             "Záměnnost",
		Incompatible:            "Neslučitelnost",
		Classes:                 "Třída(y)",
		Classification:          "Klasifikace",
		Guarantors:              "Garant(i)",
		Teachers:                "Vyučující",
		CourseGlobalNotes:       "Poznámky k předmětu",
		NotRated:                "Nehodnoceno",
		NoRatings:               "Žádná hodnocení",
		AdditionalRatings:       "Rozšířené hodnocení",
		CategoricalRatings:      "Kategorická hodnocení",
		SISLink:                 "Odkaz do SIS",
		Description:             "Popis",
		NoDescription:           "Pro tento předmět není k dispozici žádný popis.",
		Survey:                  "Anketa",
		NoSurvey:                "Nejsou k dispozici žádné komentáře.",
		SurveySearchPlaceholder: "Hledat v anketě...",
		ScrollToSearch:          "Nahoru na vyhledávání",
		Detail:                  "Detailní informace",
		NoDetail:                "Pro tento předmět nejsou k dispozici žádné podrobnosti.",
		CourseWithCode:          "Předmět s kódem ",
		NotFound:                " nenalezen.",
		// utils
		Utils: utils.Texts[language.CS],
	},
	language.EN: {
		Faculty:                 "Faculty",
		Department:              "Department",
		Semester:                "Semester",
		Winter:                  "Winter",
		Summer:                  "Summer",
		Both:                    "Both",
		WinterAssign:            "Winter",
		SummerAssign:            "Summer",
		BothAssign:              "Both",
		Unassigned:              "Unassigned",
		Year:                    "Year",
		Assign:                  "Assign",
		ECredits:                "E-Credits",
		InDegreePlan:            "In degree plan",
		Completed:               "Completed",
		Capacity:                "Capacity",
		CapacityNoLimit:         "No limit",
		ScopeExam:               "Scope, examination",
		StateOfCourse:           "State of the course",
		LanguageOfCourse:        "Language",
		AdditionalInfo:          "Additional information",
		Prerequisites:           "Pre-requisites",
		Corequisites:            "Co-requisites",
		Interchange:             "Interchangeability",
		Incompatible:            "Incompatibility",
		Classes:                 "Class(es)",
		Classification:          "Classification",
		Guarantors:              "Guarantor(s)",
		Teachers:                "Teacher(s)",
		CourseGlobalNotes:       "Course notes",
		NotRated:                "Not rated",
		NoRatings:               "No ratings",
		AdditionalRatings:       "Additional ratings",
		CategoricalRatings:      "Categorical ratings",
		SISLink:                 "Link to SIS",
		Description:             "Description",
		NoDescription:           "No description available for this course.",
		Survey:                  "Survey",
		NoSurvey:                "No comments available.",
		SurveySearchPlaceholder: "Search in survey...",
		ScrollToSearch:          "Scroll to search",
		Detail:                  "Detailed information",
		NoDetail:                "No detailed information available for this course.",
		CourseWithCode:          "Course with code ",
		NotFound:                " not found.",
		// utils
		Utils: utils.Texts[language.EN],
	},
}
