package degreeplan

import (
	"strconv"

	"github.com/michalhercik/RecSIS/language"
)

type text struct {
	pageTitle             string
	showDegreePlan        string
	chooseDegreePlan      string
	degreePlanPlaceholder string
	enrollmentYear        string
	noDegreePlanResults   string
	code                  string
	title                 string
	status                string
	completed             string
	inBlueprint           string
	unassigned            string
	notCompleted          string
	credits               string
	creditsShort          string
	needed                string
	winter                string
	summer                string
	both                  string
	guarantors            string
	blueprint             string
	assign                string
	courseIsUnassigned    string
	year                  string
	winterAssign          string
	summerAssign          string
	// language
	language language.Language
	// errors
	errInvalidDPYear   string
	errCannotGetUserDP string
	errCannotGetDP     string
	errDPNotFound      string
	errFailedDPSearch  string
	errPageNotFound    string
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
		pageTitle:             "Studijní plán",
		showDegreePlan:        "Zobrazte si studijní plán",
		chooseDegreePlan:      "Vyberte si studijní plán",
		degreePlanPlaceholder: "např. NISD23N",
		enrollmentYear:        "Rok zápisu",
		noDegreePlanResults:   "Žádný studijní plán nenalezen.",
		code:                  "Kód",
		title:                 "Název",
		status:                "Stav",
		completed:             "Splněno",
		inBlueprint:           "Blueprint",
		unassigned:            "Nezařazen",
		notCompleted:          "Nesplněno",
		credits:               "Kredity",
		creditsShort:          "Kr.",
		needed:                "potřeba",
		winter:                "ZS",
		summer:                "LS",
		both:                  "Oba",
		guarantors:            "Garant(i)",
		blueprint:             "Blueprint",
		assign:                "Přiřadit",
		courseIsUnassigned:    "Kurz je v blueprintu, ale není zařazen.",
		year:                  "ročník",
		winterAssign:          "ZS",
		summerAssign:          "LS",
		// language
		language: language.CS,
		// errors
		errInvalidDPYear:   "Neplatný rok studijního plánu",
		errCannotGetUserDP: "Nebylo možné získat studijní plán uživatele",
		errCannotGetDP:     "Nebylo možné získat vybraný studijní plán",
		errDPNotFound:      "Studijní plán nenalezen",
		errFailedDPSearch:  "Nebylo možné vyhledat studijní plány",
		errPageNotFound:    "Stránka nenalezena",
	},
	language.EN: {
		pageTitle:             "Degree Plan",
		showDegreePlan:        "Show your degree plan",
		chooseDegreePlan:      "Choose a degree plan",
		degreePlanPlaceholder: "e.g. NISD23N",
		enrollmentYear:        "Enrollment year",
		noDegreePlanResults:   "No degree plan found.",
		code:                  "Code",
		title:                 "Title",
		status:                "Status",
		completed:             "Completed",
		inBlueprint:           "Blueprint",
		unassigned:            "Unassigned",
		notCompleted:          "Not completed",
		credits:               "Credits",
		creditsShort:          "Cr.",
		needed:                "needed",
		winter:                "Winter",
		both:                  "Both",
		summer:                "Summer",
		guarantors:            "Guarantor(s)",
		blueprint:             "Blueprint",
		assign:                "Assign",
		courseIsUnassigned:    "Course is in the blueprint but not assigned.",
		year:                  "Year",
		winterAssign:          "Winter",
		summerAssign:          "Summer",
		// language
		language: language.EN,
		// errors
		errInvalidDPYear:   "Invalid degree plan year",
		errCannotGetUserDP: "Unable to retrieve user degree plan",
		errCannotGetDP:     "Unable to retrieve selected degree plan",
		errDPNotFound:      "Degree plan not found",
		errFailedDPSearch:  "Failed to search for degree plans",
		errPageNotFound:    "Page not found",
	},
}
