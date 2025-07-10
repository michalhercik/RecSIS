package degreeplan

import (
	"strconv"

	"github.com/michalhercik/RecSIS/language"
)

type text struct {
	pageTitle             string
	showDegreePlan        string
	showDegreePlanShort   string
	saveDegreePlan        string
	saveDegreePlanShort   string
	chooseDegreePlan      string
	chooseDegreePlanHelp  string
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
	language              language.Language
	errInvalidDPYear      string
	errCannotGetUserDP    string
	errCannotGetDP        string
	errCannotSaveDP       string
	errDPNotFound         string
	errFailedDPSearch     string
	errPageNotFound       string
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
		showDegreePlan:        "Vybrat studijní plán",
		showDegreePlanShort:   "Vybrat SP",
		saveDegreePlan:        "Uložit studijní plán",
		saveDegreePlanShort:   "Uložit SP",
		chooseDegreePlan:      "Vyberte si studijní plán",
		chooseDegreePlanHelp:  "Do vyhledávacího pole vlevo zadejte kód Vašeho studijního plánu a vpravo vyberte Váš rok zápisu do daného studia. Kód studijního plánu naleznete v SIS v záložce 'Osobní údaje a nastavení' v položce 'Studijní plán' za názvem v závorce. Rok zápisu naleznete taktéž v SIS v záložce 'Osobní údaje a nastavení' v položce 'Datum zápisu'.",
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
		language:              language.CS,
		errInvalidDPYear:      "Neplatný rok studijního plánu",
		errCannotGetUserDP:    "Nebylo možné získat studijní plán uživatele",
		errCannotGetDP:        "Nebylo možné získat vybraný studijní plán",
		errCannotSaveDP:       "Nebylo možné uložit studijní plán",
		errDPNotFound:         "Studijní plán nenalezen",
		errFailedDPSearch:     "Nebylo možné vyhledat studijní plány",
		errPageNotFound:       "Stránka nenalezena",
	},
	language.EN: {
		pageTitle:             "Degree Plan",
		showDegreePlan:        "Choose degree plan",
		showDegreePlanShort:   "Choose DP",
		saveDegreePlan:        "Save degree plan",
		saveDegreePlanShort:   "Save DP",
		chooseDegreePlan:      "Choose a degree plan",
		chooseDegreePlanHelp:  "Enter your degree plan code in the search field on the left and select your enrollment year on the right. You can find the degree plan code in SIS under 'Personal data and settings' tab in the 'Curriculum' item, next to the name, in parentheses. The enrollment year can also be found in SIS under 'Personal data and settings' tab in the 'Enrollment date' item.",
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
		language:              language.EN,
		errInvalidDPYear:      "Invalid degree plan year",
		errCannotGetUserDP:    "Unable to retrieve user degree plan",
		errCannotGetDP:        "Unable to retrieve selected degree plan",
		errDPNotFound:         "Degree plan not found",
		errFailedDPSearch:     "Failed to search for degree plans",
		errPageNotFound:       "Page not found",
	},
}
