package page

import (
	"github.com/michalhercik/RecSIS/language"
)

type text struct {
	logout                 string
	contact                string
	noCoursesFound         string
	searchPlaceholder      string
	searchButton           string
	language               language.Language
	errUnsupportedLanguage string
	errQuickSearchFailed   string
}

var texts = map[language.Language]text{
	language.CS: {
		logout:                 "Odhlásit se",
		contact:                "V případě jakýchkoliv problémů kontaktujte tým RecSIS.",
		noCoursesFound:         "Žádné předměty nebyly nalezeny.",
		searchPlaceholder:      "Hledej předmět...",
		searchButton:           "Hledej",
		language:               language.CS,
		errUnsupportedLanguage: "Tato jazyková verze není podporována.",
		errQuickSearchFailed:   "Chyba při rychlém vyhledávání předmětů.",
	},
	language.EN: {
		logout:                 "Logout",
		contact:                "In case of any problems, please contact the RecSIS team.",
		noCoursesFound:         "No courses found.",
		searchPlaceholder:      "Search for course...",
		searchButton:           "Search",
		language:               language.EN,
		errUnsupportedLanguage: "This language version is not supported.",
		errQuickSearchFailed:   "Error during course quick-search.",
	},
}
