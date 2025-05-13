package cas

import "github.com/michalhercik/RecSIS/language"

type Text struct {
	Title            string
	PrimaryMessage   string
	SecondaryMessage string
	LogInButton      string
}

var texts = map[language.Language]Text{
	language.CS: {
		Title:            "Odhlášení - RecSIS",
		PrimaryMessage:   "Úspešné odhlášení",
		SecondaryMessage: "Byli jste úspěšně odhlášeni.",
		LogInButton:      "Přihlásit se",
	},
	language.EN: {
		Title:            "Logout - RecSIS",
		PrimaryMessage:   "Logout Successful",
		SecondaryMessage: "You have been logged out successfully.",
		LogInButton:      "Log In",
	},
}
