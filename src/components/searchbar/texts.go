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
	// errors
	errUnsupportedLanguage string
	errQuickSearchFailed   string
}

var texts = map[language.Language]text{
	language.CS: {
		noCoursesFound:    "Žádné předměty nebyly nalezeny.",
		searchPlaceholder: "Hledej předmět...",
		searchButton:      "Hledej",
		// language
		language: language.CS,
		// errors
		errUnsupportedLanguage: "Tato jazyková verze není podporována.",
		errQuickSearchFailed:   "Chyba při rychlém vyhledávání předmětů.",
	},
	language.EN: {
		noCoursesFound:    "No courses found.",
		searchPlaceholder: "Search for course...",
		searchButton:      "Search",
		// language
		language: language.EN,
		// errors
		errUnsupportedLanguage: "This language version is not supported.",
		errQuickSearchFailed:   "Error during course quick-search.",
	},
}
