package cas

import "github.com/michalhercik/RecSIS/language"

type text struct {
	title            string
	primaryMessage   string
	secondaryMessage string
	logInButton      string
}

var texts = map[language.Language]text{
	language.CS: {
		title:            "Odhlášení - RecSIS",
		primaryMessage:   "Úspešné odhlášení",
		secondaryMessage: "Byli jste úspěšně odhlášeni.",
		logInButton:      "Přihlásit se",
	},
	language.EN: {
		title:            "Logout - RecSIS",
		primaryMessage:   "Logout Successful",
		secondaryMessage: "You have been logged out successfully.",
		logInButton:      "Log In",
	},
}
