package ui

import (
	"github.com/michalhercik/RecSIS/language"
)

type Translated struct {
	ViewAll string
	Credits string
	ViewMore string
}

var i18n = map[language.Language]Translated{
	language.CS: {
		ViewAll: "Zobrazit vše",
		Credits: "Kr.",
		ViewMore: "Zobrazit více",
	},
	language.EN: {
		ViewAll: "View All",
		Credits: "Cr.",
		ViewMore: "View More",
	},
}
