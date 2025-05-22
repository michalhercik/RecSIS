package page

import (
	"github.com/michalhercik/RecSIS/language"
)

type Text struct {
	Contact string
}

var texts = map[language.Language]Text{
	language.CS: {
		Contact: "V případě jakýchkoliv problémů kontaktujte tým RecSIS.",
	},
	language.EN: {
		Contact: "In case of any problems, please contact the RecSIS team.",
	},
}
