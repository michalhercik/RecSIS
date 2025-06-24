package page

import (
	"github.com/michalhercik/RecSIS/language"
)

type text struct {
	logout  string
	contact string
}

var texts = map[language.Language]text{
	language.CS: {
		logout:  "Odhlásit se",
		contact: "V případě jakýchkoliv problémů kontaktujte tým RecSIS.",
	},
	language.EN: {
		logout:  "Logout",
		contact: "In case of any problems, please contact the RecSIS team.",
	},
}
