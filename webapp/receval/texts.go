package receval

import (
	"github.com/michalhercik/RecSIS/language"
)

type text struct {
	pageTitle string
	winter                    string
	summer                    string
	both                      string
	credits                   string
	noGuarantors              string
	language                  language.Language
	errPageNotFound           string
}

var texts = map[language.Language]text{
	language.CS: {
		pageTitle: 				   "Evaluace",
		winter:                    "ZS",
		summer:                    "LS",
		both:                      "Oba",
		credits:                   "Kredity",
		noGuarantors:              "Žádní garanti",
		language:                  language.CS,
		errPageNotFound:           "Stránka nenalezena",
	},
	language.EN: {
		pageTitle: 				   "Evaluation",
		winter:                    "Winter",
		summer:                    "Summer",
		both:                      "Both",
		credits:                   "Credits",
		noGuarantors:              "No guarantors",
		language:                  language.EN,
		errPageNotFound:           "Page not found",
	},
}
