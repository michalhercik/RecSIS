package coursedetail

import (
	"strconv"

	"github.com/michalhercik/RecSIS/language"
	"github.com/michalhercik/RecSIS/utils"
)

type text struct {
	faculty                 string
	department              string
	semester                string
	winter                  string
	summer                  string
	both                    string
	winterAssign            string
	summerAssign            string
	bothAssign              string
	unassigned              string
	year                    string
	assign                  string
	eCredits                string
	inDegreePlan            string
	completed               string
	capacity                string
	scopeExam               string
	stateOfCourse           string
	languageOfCourse        string
	additionalInfo          string
	prerequisites           string
	corequisites            string
	interchange             string
	incompatible            string
	classes                 string
	classification          string
	guarantors              string
	teachers                string
	courseGlobalNotes       string
	notRated                string
	noRatings               string
	additionalRatings       string
	categoricalRatings      string
	sisLink                 string
	description             string
	noDescription           string
	survey                  string
	noSurvey                string
	surveySearchPlaceholder string
	cancelFilters           string
	scrollToSearch          string
	detail                  string
	noDetail                string
	courseWithCode          string
	notFound                string
	// utils
	utils utils.Text
}

func (t text) yearStr(year int) string {
	if t.utils.Language == language.CS {
		return strconv.Itoa(year) + ". " + t.year
	} else if t.utils.Language == language.EN {
		return t.year + " " + strconv.Itoa(year)
	}
	return ""
}

var texts = map[language.Language]text{
	language.CS: {
		faculty:                 "Fakulta",
		department:              "Katedra",
		semester:                "Semestr",
		winter:                  "Zimní",
		summer:                  "Letní",
		both:                    "Oba",
		winterAssign:            "ZS",
		summerAssign:            "LS",
		bothAssign:              "Oba",
		unassigned:              "Nezařazen",
		year:                    "ročník",
		assign:                  "Přiřadit",
		eCredits:                "E-Kredity",
		inDegreePlan:            "V studijním plánu",
		completed:               "Splněno",
		capacity:                "Kapacita",
		scopeExam:               "Rozsah, examinace",
		stateOfCourse:           "Stav předmětu",
		languageOfCourse:        "Jazyk výuky",
		additionalInfo:          "Další informace",
		prerequisites:           "Prerekvizity",
		corequisites:            "Korekvizity",
		interchange:             "Záměnnost",
		incompatible:            "Neslučitelnost",
		classes:                 "Třída(y)",
		classification:          "Klasifikace",
		guarantors:              "Garant(i)",
		teachers:                "Vyučující",
		courseGlobalNotes:       "Poznámky k předmětu",
		notRated:                "Nehodnoceno",
		noRatings:               "Žádná hodnocení",
		additionalRatings:       "Rozšířené hodnocení",
		categoricalRatings:      "Kategorická hodnocení",
		sisLink:                 "Odkaz do SIS",
		description:             "Popis",
		noDescription:           "Pro tento předmět není k dispozici žádný popis.",
		survey:                  "Anketa",
		noSurvey:                "Nejsou k dispozici žádné komentáře.",
		surveySearchPlaceholder: "Hledat v anketě...",
		cancelFilters:           "Zrušit vybrané filtry",
		scrollToSearch:          "Nahoru na vyhledávání",
		detail:                  "Detailní informace",
		noDetail:                "Pro tento předmět nejsou k dispozici žádné podrobnosti.",
		courseWithCode:          "Předmět s kódem ",
		notFound:                " nenalezen.",
		// utils
		utils: utils.Texts[language.CS],
	},
	language.EN: {
		faculty:                 "Faculty",
		department:              "Department",
		semester:                "Semester",
		winter:                  "Winter",
		summer:                  "Summer",
		both:                    "Both",
		winterAssign:            "Winter",
		summerAssign:            "Summer",
		bothAssign:              "Both",
		unassigned:              "Unassigned",
		year:                    "Year",
		assign:                  "Assign",
		eCredits:                "E-Credits",
		inDegreePlan:            "In degree plan",
		completed:               "Completed",
		capacity:                "Capacity",
		scopeExam:               "Scope, examination",
		stateOfCourse:           "State of the course",
		languageOfCourse:        "Language",
		additionalInfo:          "Additional information",
		prerequisites:           "Pre-requisites",
		corequisites:            "Co-requisites",
		interchange:             "Interchangeability",
		incompatible:            "Incompatibility",
		classes:                 "Class(es)",
		classification:          "Classification",
		guarantors:              "Guarantor(s)",
		teachers:                "Teacher(s)",
		courseGlobalNotes:       "Course notes",
		notRated:                "Not rated",
		noRatings:               "No ratings",
		additionalRatings:       "Additional ratings",
		categoricalRatings:      "Categorical ratings",
		sisLink:                 "Link to SIS",
		description:             "Description",
		noDescription:           "No description available for this course.",
		survey:                  "Survey",
		noSurvey:                "No comments available.",
		surveySearchPlaceholder: "Search in survey...",
		cancelFilters:           "Cancel selected filters",
		scrollToSearch:          "Scroll to search",
		detail:                  "Detailed information",
		noDetail:                "No detailed information available for this course.",
		courseWithCode:          "Course with code ",
		notFound:                " not found.",
		// utils
		utils: utils.Texts[language.EN],
	},
}
