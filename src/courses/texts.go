package courses

import (
	"github.com/michalhercik/RecSIS/language"
)

type text struct {
	pageTitle         string
	winter            string
	summer            string
	both              string
	w                 string
	s                 string
	n                 string
	er                string
	un                string
	searchPlaceholder string
	searchButton      string
	filterButton      string
	showResults       string
	topFilter         string
	showMore4Minus    string
	showMore5Plus     string
	showLess          string
	cancelFilters     string
	assign            string
	year              string
	winterAssign      string
	summerAssign      string
	credits           string
	noGuarantors      string
	readMore          string
	readLess          string
	inDegreePlan      string
	loadMore          string
	previousPage      string
	nextPage          string
	page              string
	of                string
	noCoursesFound    string
	// language
	language language.Language
	// errors
	errUnexpectedNumberOfCourses string
	errCannotLoadCourses         string
	errCannotSearchCourses       string
	errInvalidPageNumber         string
	errInvalidNumberOfCourses    string
	errPageNotFound              string
}

func (t text) showMore(rest int) string {
	if rest < 5 {
		return t.showMore4Minus
	} else {
		return t.showMore5Plus
	}
}

var texts = map[language.Language]text{
	language.CS: {
		pageTitle:         "Hledání",
		winter:            "Zimní",
		summer:            "Letní",
		both:              "Oba",
		w:                 "ZS",
		s:                 "LS",
		n:                 "ER",
		er:                "ER",
		un:                "NEZAŘ.",
		searchPlaceholder: "Hledej...",
		searchButton:      "Hledej",
		filterButton:      "Zobrazit filtrování",
		showResults:       "Zobrazit výsledky",
		topFilter:         "Filtrovat",
		showMore4Minus:    "Další",
		showMore5Plus:     "Dalších",
		showLess:          "Skrýt ostatní",
		cancelFilters:     "Zrušit vybrané filtry",
		assign:            "Přiřadit",
		year:              "ročník",
		winterAssign:      "ZS",
		summerAssign:      "LS",
		credits:           "Kredity",
		noGuarantors:      "Žádní garanti",
		readMore:          "Číst dále",
		readLess:          "Sbalit",
		inDegreePlan:      "Studijní plán",
		loadMore:          "Zobrazit další",
		previousPage:      "Předchozí stránka",
		nextPage:          "Další stránka",
		page:              "Strana",
		of:                "z",
		noCoursesFound:    "Žádné předměty nebyly nalezeny.",
		// language
		language: language.CS,
		// errors
		errUnexpectedNumberOfCourses: "Neočekávaný počet předmětů pro přiřazení do Blueprintu",
		errCannotLoadCourses:         "Nelze načíst předměty z databáze",
		errCannotSearchCourses:       "Nelze vyhledat předměty",
		errInvalidPageNumber:         "neplatné číslo stránky: musí být celé kladné číslo",
		errInvalidNumberOfCourses:    "neplatný počet předmětů na stránku: musí být celé kladné číslo",
		errPageNotFound:              "Stránka nenalezena",
	},
	language.EN: {
		pageTitle:         "Search",
		winter:            "Winter",
		summer:            "Summer",
		both:              "Both",
		w:                 "WIN",
		s:                 "SUM",
		n:                 "ER",
		er:                "ER",
		un:                "UNASS.",
		searchPlaceholder: "Search...",
		searchButton:      "Search",
		filterButton:      "Show filters",
		showResults:       "Show results",
		topFilter:         "Filter",
		showMore4Minus:    "Another",
		showMore5Plus:     "Another",
		showLess:          "Hide others",
		cancelFilters:     "Cancel selected filters",
		assign:            "Assign",
		year:              "Year",
		winterAssign:      "Winter",
		summerAssign:      "Summer",
		credits:           "Credits",
		noGuarantors:      "No guarantors",
		readMore:          "Read more",
		readLess:          "Hide",
		inDegreePlan:      "Degree plan",
		loadMore:          "Show more",
		previousPage:      "Previous page",
		nextPage:          "Next page",
		page:              "Page",
		of:                "of",
		noCoursesFound:    "No courses found.",
		// language
		language: language.EN,
		// errors
		errUnexpectedNumberOfCourses: "Unexpected number of courses for Blueprint assignment",
		errCannotLoadCourses:         "Cannot load courses from the database",
		errCannotSearchCourses:       "Cannot search for courses",
		errInvalidPageNumber:         "Invalid page number: must be a whole positive number",
		errInvalidNumberOfCourses:    "Invalid number of courses per page: must be a whole positive number",
		errPageNotFound:              "Page not found",
	},
}
