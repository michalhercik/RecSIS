package cas

import "github.com/michalhercik/RecSIS/language"

type text struct {
	title                 string
	logoutSuccessHeadline string
	logoutSuccessMessage  string
	loginButton           string
	// errors
	errUnauthorized string
}

var texts = map[language.Language]text{
	language.CS: {
		title:                 "Odhlášení - RecSIS",
		logoutSuccessHeadline: "Úspěšné odhlášení",
		logoutSuccessMessage:  "Byli jste úspěšně odhlášeni.",
		loginButton:           "Přihlásit se",
		// errors
		errUnauthorized: "Neoprávněný přístup",
	},
	language.EN: {
		title:                 "Logout - RecSIS",
		logoutSuccessHeadline: "Logout Successful",
		logoutSuccessMessage:  "You have been logged out successfully.",
		loginButton:           "Log In",
		// errors
		errUnauthorized: "Unauthorized Access",
	},
}
