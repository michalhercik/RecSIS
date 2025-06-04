package utils

import (
	"github.com/michalhercik/RecSIS/language"
)

type Text struct {
	Language   language.Language
	Home       string
	Courses    string
	Blueprint  string
	DegreePlan string
	Login      string
	Contact    string
}

var Texts = map[language.Language]Text{
	language.CS: {
		Language:   language.CS,
		Home:       "Domů",
		Courses:    "Hledání",
		Blueprint:  "Blueprint",
		DegreePlan: "Studijní plán",
		Login:      "Přihlášení",
		Contact:    "V případě jakýchkoliv problémů kontaktujte tým RecSIS.",
	},
	language.EN: {
		Language:   language.EN,
		Home:       "Home",
		Courses:    "Search",
		Blueprint:  "Blueprint",
		DegreePlan: "Degree plan",
		Login:      "Login",
		Contact:    "In case of any problems, please contact the RecSIS team.",
	},
}
