package searchbar

import (
	"github.com/michalhercik/RecSIS/language"
)

type text struct {
	noCoursesFound    string
	searchPlaceholder string
	searchButton      string
	// language
	language language.Language
}

var texts = map[language.Language]text{
	language.CS: {
		noCoursesFound:    "Žádné předměty nebyly nalezeny.",
		searchPlaceholder: "Hledej předmět...",
		searchButton:      "Hledej",
		// language
		language: language.CS,
	},
	language.EN: {
		noCoursesFound:    "No courses found.",
		searchPlaceholder: "Search for course...",
		searchButton:      "Search",
		// language
		language: language.EN,
	},
}
