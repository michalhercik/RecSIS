package searchbar

import (
	"github.com/michalhercik/RecSIS/language"
	"github.com/michalhercik/RecSIS/utils"
)

type Text struct {
	NoCoursesFound    string
	SearchPlaceholder string
	SearchButton      string
	Utils             utils.Text
}

var texts = map[language.Language]Text{
	language.CS: {
		NoCoursesFound:    "Žádné předměty nebyly nalezeny.",
		SearchPlaceholder: "Hledej předmět...",
		SearchButton:      "Hledej",
		Utils:             utils.Texts[language.CS],
	},
	language.EN: {
		NoCoursesFound:    "No courses found.",
		SearchPlaceholder: "Search for course...",
		SearchButton:      "Search",
		Utils:             utils.Texts[language.EN],
	},
}
