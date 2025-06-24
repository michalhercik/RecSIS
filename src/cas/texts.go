package cas

import "github.com/michalhercik/RecSIS/language"

type text struct {
	title                 string
	logoutSuccessHeadline string
	logoutSuccessMessage  string
	loginButton           string
	// errors
	errUnauthorized               string
	errCannotGetUserIDFromSession string
	errCannotCreateSession        string
	errCannotLogout               string
	errCannotCreateUser           string
	errCannotGetTicket            string
}

var texts = map[language.Language]text{
	language.CS: {
		title:                 "Odhlášení - RecSIS",
		logoutSuccessHeadline: "Úspěšné odhlášení",
		logoutSuccessMessage:  "Byli jste úspěšně odhlášeni.",
		loginButton:           "Přihlásit se",
		// errors
		errUnauthorized:               "Neoprávněný přístup",
		errCannotGetUserIDFromSession: "Nepodařilo se získat ID uživatele ze session.",
		errCannotCreateSession:        "Nepodařilo se vytvořit session.",
		errCannotLogout:               "Nepodařilo se odhlásit.",
		errCannotCreateUser:           "Nepodařilo se vytvořit uživatele.",
		errCannotGetTicket:            "Nepodařilo se získat ticket z požadavku.",
	},
	language.EN: {
		title:                 "Logout - RecSIS",
		logoutSuccessHeadline: "Logout Successful",
		logoutSuccessMessage:  "You have been logged out successfully.",
		loginButton:           "Log In",
		// errors
		errUnauthorized:               "Unauthorized Access",
		errCannotGetUserIDFromSession: "Failed to retrieve user ID from session.",
		errCannotCreateSession:        "Failed to create session.",
		errCannotLogout:               "Failed to log out.",
		errCannotCreateUser:           "Failed to create user.",
		errCannotGetTicket:            "Failed to retrieve ticket from request.",
	},
}
