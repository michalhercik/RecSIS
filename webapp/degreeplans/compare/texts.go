package compare

import (
	"github.com/michalhercik/RecSIS/language"
)

type text struct {
	pageTitle         string
	switchPlans       string
	showDifferences   string
	showSame          string
	credits           string
	locateInOtherPlan string
	// language
	language language.Language
	// errors
	errCannotGetDP   string
	errDPCodeMissing string
	errDPNotExisting string
	errPageNotFound  string
}

var texts = map[language.Language]text{
	language.CS: {
		pageTitle:         "Porovnání plánů",
		switchPlans:       "Prohodit plány",
		showDifferences:   "Zobrazit rozdíly",
		showSame:          "Zobrazit shodné předměty",
		credits:           "Kredity",
		locateInOtherPlan: "Najít v druhém plánu",
		// language
		language: language.CS,
		// errors
		errCannotGetDP:   "Nepodařilo se získat studijní plán z databáze.",
		errDPCodeMissing: "Chybí kód studijního plánu",
		errDPNotExisting: "Zadaný studijní plán neexistuje.",
		errPageNotFound:  "Stránka nebyla nalezena.",
	},
	language.EN: {
		pageTitle:         "Compare Plans",
		switchPlans:       "Switch Plans",
		showDifferences:   "Show Differences",
		showSame:          "Show Same Courses",
		credits:           "Credits",
		locateInOtherPlan: "Find in the other plan",
		// language
		language: language.EN,
		// errors
		errCannotGetDP:   "Failed to get degree plan from the database.",
		errDPCodeMissing: "Degree plan code is missing.",
		errDPNotExisting: "The selected degree plan does not exist.",
		errPageNotFound:  "Page not found.",
	},
}
