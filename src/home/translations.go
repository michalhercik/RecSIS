package home

import (
	"github.com/michalhercik/RecSIS/language"
	"github.com/michalhercik/RecSIS/utils"
)

type text struct {
	PageTitle          string
	Welcome            string
	RecSISIntro        string
	RecommendedCourses string
	NewCourses         string
	Winter             string
	Summer             string
	Both               string
	Credits            string
	NoGuarantors       string
	// utils
	Utils utils.Text
}

var texts = map[language.Language]text{
	language.CS: {
		PageTitle:          "Domů",
		Welcome:            "Vítejte!",
		RecSISIntro:        "RecSIS je systém pro plánování studia, kontrolování studijních povinností a doporučování kurzů.",
		RecommendedCourses: "Doporučené kurzy přímo pro vás",
		NewCourses:         "Nové kurzy",
		Winter:             "ZS",
		Summer:             "LS",
		Both:               "Oba",
		Credits:            "Kredity",
		NoGuarantors:       "Žádní garanti",
		// utils
		Utils: utils.Texts[language.CS],
	},
	language.EN: {
		PageTitle:          "Home",
		Welcome:            "Welcome!",
		RecSISIntro:        "RecSIS is a system for study planning, monitoring study obligations, and recommending courses.",
		RecommendedCourses: "Recommended courses just for you",
		NewCourses:         "New courses",
		Winter:             "Winter",
		Summer:             "Summer",
		Both:               "Both",
		Credits:            "Credits",
		NoGuarantors:       "No guarantors",
		// utils
		Utils: utils.Texts[language.EN],
	},
}
