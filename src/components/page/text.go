package page

import (
	"github.com/michalhercik/RecSIS/language"
	"github.com/michalhercik/RecSIS/utils"
)

type Text struct {
	Contact string
	Utils   utils.Text
}

var texts = map[language.Language]Text{
	language.CS: {
		Contact: "V případě jakýchkoliv problémů kontaktujte tým RecSIS.",
	},
	language.EN: {
		Contact: "In case of any problems, please contact the RecSIS team.",
	},
}
