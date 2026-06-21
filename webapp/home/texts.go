package home

import (
	"github.com/michalhercik/RecSIS/language"
)

type text struct {
	pageTitle                 string
	recommendedCourses        string
	newCourses                string
	winter                    string
	summer                    string
	both                      string
	credits                   string
	noGuarantors              string
	language                  language.Language
	errRecommenderUnavailable string
	errCannotLoadCourses      string
	errPageNotFound           string
	errForYou                 string
}

var texts = map[language.Language]text{
	language.CS: {
		pageTitle:                 "Domů",
		recommendedCourses:        "Pro Tebe",
		newCourses:                "Nové kurzy",
		winter:                    "ZS",
		summer:                    "LS",
		both:                      "Oba",
		credits:                   "Kredity",
		noGuarantors:              "Žádní garanti",
		language:                  language.CS,
		errRecommenderUnavailable: "Nelze se připojit k doporučovacímu systému",
		errCannotLoadCourses:      "Nelze načíst kurzy na stránce",
		errPageNotFound:           "Stránka nenalezena",
		errForYou:                 "Nelze načíst ForYou doporučení",
	},
	language.EN: {
		pageTitle:                 "Home",
		recommendedCourses:        "For You",
		newCourses:                "New courses",
		winter:                    "Winter",
		summer:                    "Summer",
		both:                      "Both",
		credits:                   "Credits",
		noGuarantors:              "No guarantors",
		language:                  language.EN,
		errRecommenderUnavailable: "Cannot connect to recommender system",
		errCannotLoadCourses:      "Cannot load courses on the page",
		errPageNotFound:           "Page not found",
		errForYou:                 "Cannot load ForYou recommendations",
	},
}
