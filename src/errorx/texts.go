package errorx

import (
	"github.com/michalhercik/RecSIS/language"
)

type text struct {
	errorText string
	// language
	language language.Language
	// errors
	errOk                    string
	errGeneric               string
	errCannotRenderPage      string
	errCannotRenderComponent string
}

var texts = map[language.Language]text{
	language.CS: {
		errorText: "Error",
		// language
		language: language.CS,
		// errors
		errOk:                    "Vše by mělo být v pořádku",
		errGeneric:               "Došlo k chybě, zkuste to prosím znovu později.",
		errCannotRenderPage:      "Stránku se nepodařilo vygenerovat, zkuste to prosím znovu později.",
		errCannotRenderComponent: "Komponentu se nepodařilo vygenerovat, zkuste to prosím znovu později.",
	},
	language.EN: {
		errorText: "Error",
		// language
		language: language.EN,
		// errors
		errOk:                    "Everything should be fine",
		errGeneric:               "An error occurred, please try again later.",
		errCannotRenderPage:      "The page could not be rendered, please try again later.",
		errCannotRenderComponent: "The component could not be rendered, please try again later.",
	},
}
