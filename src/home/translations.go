package home

import (
	"github.com/michalhercik/RecSIS/language"
	"github.com/michalhercik/RecSIS/utils"
)

type text struct {
	PageTitle        string
	Language         string
	Introduction     string
	IntroductionText string
	HowToUse         string
	HowToUseText     string
	Authors          string
	AuthorsText      string
	// utils
	Utils utils.Text
}

var texts = map[language.Language]text{
	language.CS: {
		PageTitle:        "Domů",
		Language:         "cs",
		Introduction:     "Úvod",
		IntroductionText: "Příliš žluťoučký kůň úpěl ďábelské ódy. Nechť již hříšné saxofony ďáblů rozezvučí síň úděsnými tóny waltzu, polky a quickstepu.",
		HowToUse:         "Jak používat",
		HowToUseText:     "Používejte! Příliš žluťoučký kůň úpěl ďábelské ódy. Nechť již hříšné saxofony ďáblů rozezvučí síň úděsnými tóny waltzu, polky a quickstepu.",
		Authors:          "Autoři",
		AuthorsText:      "Jeho Milost, svobodný pán z Malé Strany, kancléř univerzitní rady, doc. Mgr. <b>Michal Hercík</b>, Th.D., LL.M., kustod historických rukopisů, hlavní kronikář akademického senátu, poradce císařské rady pro vzdělanost a vědu <br> a <br> Jeho Excelence, arcibiskup pražský, rytíř Řádu sv. Václava, prof. Ing. <b>Michal Medek</b>, Ph.D., DSc., MBA, knihovník královské univerzitní sbírky, správce archivů svaté katedrály, čestný člen spolku staroměstských alchymistů",
		// utils
		Utils: utils.Texts["cs"],
	},
	language.EN: {
		PageTitle:        "Home",
		Language:         "en",
		Introduction:     "Introduction",
		IntroductionText: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		HowToUse:         "How to use",
		HowToUseText:     "Use it! Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		Authors:          "Authors",
		AuthorsText:      "RecSIS Team!",
		// utils
		Utils: utils.Texts["en"],
	},
}
