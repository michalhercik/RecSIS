package degreeplans

import (
	"github.com/michalhercik/RecSIS/language"
)

type text struct {
	pageTitle           string
	searchPlaceholder   string
	showMore4Minus      string
	showMore5Plus       string
	showLess            string
	filterButton        string
	showResults         string
	cancelFilters       string
	noDegreePlanResults string
	code                string
	title               string
	studyType           string
	validity            string
	selectForCompare    string
	unselectForCompare  string
	compareWithSelected string
	// language
	language language.Language
	// errors
	errFailedDPSearch    string
	errFailedFetchDPMeta string
	errPageNotFound      string
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
		pageTitle:           "Vyhledávání plánů",
		searchPlaceholder:   "hledejte studijní plány podle kódu nebo názvu...",
		showMore4Minus:      "Další",
		showMore5Plus:       "Dalších",
		showLess:            "Skrýt ostatní",
		filterButton:        "Zobrazit filtrování",
		showResults:         "Zobrazit výsledky",
		cancelFilters:       "Zrušit vybrané filtry",
		noDegreePlanResults: "Nebyly nalezeny žádné studijní plány.",
		code:                "Kód",
		title:               "Název",
		studyType:           "Druh",
		validity:            "Platnost",
		selectForCompare:    "Vybrat k porovnání",
		unselectForCompare:  "Zrušit výběr k porovnání",
		compareWithSelected: "Porovnat s vybraným plánem",
		// language
		language: language.CS,
		// errors
		errFailedDPSearch:    "Nebylo možné vyhledat studijní plány",
		errFailedFetchDPMeta: "Nebylo možné načíst seznam studijních plánů",
		errPageNotFound:      "Stránka nenalezena",
	},
	language.EN: {
		pageTitle:           "Search Degree Plans",
		searchPlaceholder:   "search degree plans by code or title...",
		showMore4Minus:      "Another",
		showMore5Plus:       "Another",
		showLess:            "Hide others",
		filterButton:        "Show filters",
		showResults:         "Show results",
		cancelFilters:       "Cancel selected filters",
		noDegreePlanResults: "No degree plan found.",
		code:                "Code",
		title:               "Title",
		studyType:           "Type",
		validity:            "Validity",
		selectForCompare:    "Select for compare",
		unselectForCompare:  "Unselect for compare",
		compareWithSelected: "Compare with selected plan",
		// language
		language: language.EN,
		// errors
		errFailedDPSearch:    "Failed to search for degree plans",
		errFailedFetchDPMeta: "Failed to fetch degree plan list",
		errPageNotFound:      "Page not found",
	},
}
