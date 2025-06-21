package home

import (
	"github.com/michalhercik/RecSIS/language"
)

type text struct {
	pageTitle          string
	welcome            string
	recsisIntro        string
	recommendedCourses string
	newCourses         string
	winter             string
	summer             string
	both               string
	credits            string
	noGuarantors       string
	// language
	language language.Language
	// errors
	errRecommenderUnavailable string
	errCannotLoadCourses      string
	errPageNotFound           string
}

var texts = map[language.Language]text{
	language.CS: {
		pageTitle:          "Domů",
		welcome:            "Vítejte!",
		recsisIntro:        "RecSIS je systém pro plánování studia, kontrolování studijních povinností a doporučování kurzů.",
		recommendedCourses: "Doporučené kurzy přímo pro vás",
		newCourses:         "Nové kurzy",
		winter:             "ZS",
		summer:             "LS",
		both:               "Oba",
		credits:            "Kredity",
		noGuarantors:       "Žádní garanti",
		// language
		language: language.CS,
		// errors
		errRecommenderUnavailable: "Nelze se připojit k doporučovacímu systému",
		errCannotLoadCourses:      "Nelze načíst kurzy na stránce",
		errPageNotFound:           "Stránka nenalezena",
	},
	language.EN: {
		pageTitle:          "Home",
		welcome:            "Welcome!",
		recsisIntro:        "RecSIS is a system for study planning, monitoring study obligations, and recommending courses.",
		recommendedCourses: "Recommended courses just for you",
		newCourses:         "New courses",
		winter:             "Winter",
		summer:             "Summer",
		both:               "Both",
		credits:            "Credits",
		noGuarantors:       "No guarantors",
		// language
		language: language.EN,
		// errors
		errRecommenderUnavailable: "Cannot connect to recommender system",
		errCannotLoadCourses:      "Cannot load courses on the page",
		errPageNotFound:           "Page not found",
	},
}
