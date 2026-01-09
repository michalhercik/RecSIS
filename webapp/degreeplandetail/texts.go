package degreeplandetail

import (
	"strconv"

	"github.com/michalhercik/RecSIS/language"
)

type text struct {
	pageTitle             string
	saveDegreePlan        string
	saveDegreePlanShort   string
	removeSavedDegreePlan string
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
	errCannotGetUserDP     string
	errCannotGetDP         string
	errCannotSaveDP        string
	errCannotDeleteSavedDP string
	errDPNotFound          string
	errPageNotFound        string
	// tooltips
	ttAssignedCredits  string
	ttBlueprintCredits string
	ttCompletedCredits string
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
		saveDegreePlan:        "Uložit studijní plán",
		saveDegreePlanShort:   "Uložit SP",
		removeSavedDegreePlan: "Odstranit uložený studijní plán",
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
		language:           language.CS,
		errCannotGetUserDP: "Nebylo možné získat studijní plán uživatele",
		// errors
		errCannotGetDP:         "Nebylo možné získat vybraný studijní plán",
		errCannotSaveDP:        "Nebylo možné uložit studijní plán",
		errCannotDeleteSavedDP: "Nebylo možné smazat uložený studijní plán",
		errDPNotFound:          "Studijní plán nenalezen",
		errPageNotFound:        "Stránka nenalezena",
		// tooltips
		ttAssignedCredits:  "počet kreditů přiřazených do ročníků / limit skupiny",
		ttBlueprintCredits: "počet kreditů v blueprintu / limit skupiny",
		ttCompletedCredits: "počet splněných kreditů / limit skupiny",
	},
	language.EN: {
		pageTitle:             "Degree Plan",
		saveDegreePlan:        "Save degree plan",
		saveDegreePlanShort:   "Save DP",
		removeSavedDegreePlan: "Remove saved degree plan",
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
		errCannotGetUserDP:     "Unable to retrieve user degree plan",
		errCannotGetDP:         "Unable to retrieve selected degree plan",
		errCannotSaveDP:        "Unable to save degree plan",
		errCannotDeleteSavedDP: "Unable to delete saved degree plan",
		errDPNotFound:          "Degree plan not found",
		errPageNotFound:        "Page not found",
		// tooltips
		ttAssignedCredits:  "sum of credits assigned to years / group limit",
		ttBlueprintCredits: "sum of credits in the blueprint / group limit",
		ttCompletedCredits: "sum of completed credits / group limit",
	},
}
