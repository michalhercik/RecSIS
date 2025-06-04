package page

import (
	"github.com/michalhercik/RecSIS/language"
)

type Text struct {
	Logout  string
	Contact string
}

var texts = map[language.Language]Text{
	language.CS: {
		Logout:  "Odhlásit se",
		Contact: "V případě jakýchkoliv problémů kontaktujte tým RecSIS.",
	},
	language.EN: {
		Logout:  "Logout",
		Contact: "In case of any problems, please contact the RecSIS team.",
	},
}
