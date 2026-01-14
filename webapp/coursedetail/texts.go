package coursedetail

import (
	"strconv"

	"github.com/michalhercik/RecSIS/language"
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
	and                     string
	or                      string
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
	// language
	language language.Language
	// errors
	errCourseNotFoundPre         string
	errCourseNotFoundSuf         string
	errCannotGetCourse           string
	errCannotGetCourseRatings    string
	errCannotGetRequisites       string
	errRatingMustBeInt           string
	errInvalidRating0to10        string
	errInvalidRating0or1         string
	errCannotRateCategory        string
	errCannotGetUpdatedRatings   string
	errCannotDeleteRating        string
	errUnableToRateCourse        string
	errUnexpectedNumberOfCourses string
	errCannotLoadSurvey          string
	errCannotSearchForSurvey     string
	errPageTitle                 string
	errPageNotFound              string
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
		inDegreePlan:            "Studijní plán",
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
		and:                     "a",
		or:                      "nebo",
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
		// language
		language: language.CS,
		// errors
		errCourseNotFoundPre:         "Předmět s kódem ",
		errCourseNotFoundSuf:         " nebyl nalezen",
		errCannotGetCourse:           "Nebylo možné získat předmět z databáze",
		errCannotGetCourseRatings:    "Nebylo možné získat hodnocení předmětu z databáze",
		errCannotGetRequisites:       "Nebylo možné získat rekvizity z databáze",
		errRatingMustBeInt:           "Hodnocení musí být celé číslo",
		errInvalidRating0to10:        "Hodnocení musí být v rozsahu 0-10",
		errInvalidRating0or1:         "Hodnocení musí být 0 nebo 1",
		errCannotRateCategory:        "Nebylo možné ohodnotit kategorii",
		errCannotGetUpdatedRatings:   "Nebylo možné získat aktualizovaná hodnocení",
		errCannotDeleteRating:        "Nebylo možné smazat hodnocení",
		errUnableToRateCourse:        "Nebylo možné ohodnotit předmět",
		errUnexpectedNumberOfCourses: "Neočekávaný počet předmětů pro přiřazení do Blueprintu",
		errCannotLoadSurvey:          "Nebylo možné načíst anketu",
		errCannotSearchForSurvey:     "vyhledávání selhalo",
		errPageTitle:                 "Detail předmětu",
		errPageNotFound:              "Stránka nenalezena",
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
		inDegreePlan:            "Degree plan",
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
		and:                     "and",
		or:                      "or",
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
		// language
		language: language.EN,
		// errors
		errCourseNotFoundPre:         "Course with code ",
		errCourseNotFoundSuf:         " was not found",
		errCannotGetCourse:           "Unable to retrieve course from database",
		errCannotGetCourseRatings:    "Unable to retrieve course ratings from database",
		errCannotGetRequisites:       "Unable to retrieve requisites from database",
		errRatingMustBeInt:           "Rating must be an integer",
		errInvalidRating0to10:        "Rating must be between 0 and 10",
		errInvalidRating0or1:         "Rating must be 0 or 1",
		errCannotRateCategory:        "Unable to rate category",
		errCannotGetUpdatedRatings:   "Unable to retrieve updated ratings from database",
		errCannotDeleteRating:        "Unable to delete rating",
		errUnableToRateCourse:        "Unable to rate course",
		errUnexpectedNumberOfCourses: "Unexpected number of courses for Blueprint assignment",
		errCannotLoadSurvey:          "Unable to load survey",
		errCannotSearchForSurvey:     "Search failed",
		errPageTitle:                 "Course Detail",
		errPageNotFound:              "Page not found",
	},
}
