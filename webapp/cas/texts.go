package cas

import "github.com/michalhercik/RecSIS/language"

type text struct {
	title                         string
	logoutSuccessHeadline         string
	logoutSuccessMessage          string
	loginTitle                    string
	loginMessage                  string
	loginCASMessage               string
	loginContactMessage           string
	learnMore                     string
	loginButton                   string
	errUnauthorized               string
	errCannotGetUserIDFromSession string
	errCannotCreateSession        string
	errCannotLogout               string
	errCannotCreateUser           string
	errCannotGetTicket            string
}

var texts = map[language.Language]text{
	language.CS: {
		title:                         "Odhlášení - RecSIS",
		logoutSuccessHeadline:         "Úspěšné odhlášení",
		logoutSuccessMessage:          "Byli jste úspěšně odhlášeni.",
		loginTitle:                    "RecSIS",
		loginMessage:                  "RecSIS je chytřejší a rychlejší alternativou ke kartě předmětů v SIS. Rychle vyhledávejte, filtrujte a prozkoumávejte podrobné informace o kurzech. Naplánujte si studium snadno už dnes a brzy se těšte na personalizované doporučování kurzů.",
		loginCASMessage:               "Přihlaste se do RecSIS pomocí externí autentizační služby %s. Pokud v RecSISu účet nemáte, systém se ho pokusí vytvořit.",
		loginContactMessage:           "V případě jakýchkoliv problémů kontaktujte podporu na adrese %s.",
		learnMore:                     "Zjistit více",
		loginButton:                   "Přihlásit se",
		errUnauthorized:               "Neoprávněný přístup",
		errCannotGetUserIDFromSession: "Nepodařilo se získat ID uživatele ze session.",
		errCannotCreateSession:        "Nepodařilo se vytvořit session.",
		errCannotLogout:               "Nepodařilo se odhlásit.",
		errCannotCreateUser:           "Nepodařilo se vytvořit uživatele.",
		errCannotGetTicket:            "Nepodařilo se získat ticket z požadavku.",
	},
	language.EN: {
		title:                         "Logout - RecSIS",
		logoutSuccessHeadline:         "Logout Successful",
		logoutSuccessMessage:          "You have been logged out successfully.",
		loginTitle:                    "RecSIS",
		loginMessage:                  "RecSIS is a smarter, faster alternative to the SIS subject page. Quickly search, filter, and explore detailed course info. Plan your studies with ease today, and get ready for personalized course recommendations coming soon.",
		loginCASMessage:               "Log in into RecSIS using external authentication service %s. If you do not have an account in RecSIS, it will attempt to create one.",
		loginContactMessage:           "In case of any problems, contact the support %s.",
		learnMore:                     "Learn more",
		loginButton:                   "Log In",
		errUnauthorized:               "Unauthorized Access",
		errCannotGetUserIDFromSession: "Failed to retrieve user ID from session.",
		errCannotCreateSession:        "Failed to create session.",
		errCannotLogout:               "Failed to log out.",
		errCannotCreateUser:           "Failed to create user.",
		errCannotGetTicket:            "Failed to retrieve ticket from request.",
	},
}
