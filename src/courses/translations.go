package courses

import (
	"strconv"

	"github.com/michalhercik/RecSIS/language"
	"github.com/michalhercik/RecSIS/utils"
)

type text struct {
	PageTitle         string
	Language          string
	Winter            string
	Summer            string
	Both              string
	W                 string
	S                 string
	N                 string
	ER                string
	UN                string
	SearchPlaceholder string
	SearchButton      string
	FilterButton      string
	ShowResults       string
	TopFilter         string
	ShowMore4Minus    string
	ShowMore5Plus     string
	ShowLess          string
	CancelFilters     string
	Assign            string
	Year              string
	WinterAssign      string
	SummerAssign      string
	Credits           string
	NoGuarantors      string
	ReadMore          string
	ReadLess          string
	LoadMore          string
	PreviousPage      string
	NextPage          string
	Page              string
	Of                string
	NoCoursesFound    string
	// utils
	Utils utils.Text
}

func (t text) YearStr(year int) string {
	if t.Language == "cs" {
		return strconv.Itoa(year) + ". " + t.Year
	} else if t.Language == "en" {
		return t.Year + " " + strconv.Itoa(year)
	}
	return ""
}

func (t text) ShowMore(rest int) string {
	if rest < 5 {
		return t.ShowMore4Minus
	} else {
		return t.ShowMore5Plus
	}
}

var texts = map[language.Language]text{
	language.CS: {
		PageTitle:         "Hledání",
		Language:          "cs",
		Winter:            "Zimní",
		Summer:            "Letní",
		Both:              "Oba",
		W:                 "ZS",
		S:                 "LS",
		N:                 "ER",
		ER:                "ER",
		UN:                "NEZAŘ.",
		SearchPlaceholder: "Hledej...",
		SearchButton:      "Hledej",
		FilterButton:      "Zobrazit filtrování",
		ShowResults:       "Zobrazit výsledky",
		TopFilter:         "Filtrovat",
		ShowMore4Minus:    "Další",
		ShowMore5Plus:     "Dalších",
		ShowLess:          "Skrýt ostatní",
		CancelFilters:     "Zrušit vybrané filtry",
		Assign:            "Přiřadit",
		Year:              "ročník",
		WinterAssign:      "ZS",
		SummerAssign:      "LS",
		Credits:           "Kredity",
		NoGuarantors:      "Žádní garanti",
		ReadMore:          "Číst dále",
		ReadLess:          "Sbalit",
		LoadMore:          "Zobrazit další",
		PreviousPage:      "Předchozí stránka",
		NextPage:          "Další stránka",
		Page:              "Strana",
		Of:                "z",
		NoCoursesFound:    "Žádné předměty nebyly nalezeny.",
		// utils
		Utils: utils.Texts["cs"],
	},
	language.EN: {
		PageTitle:         "Search",
		Language:          "en",
		Winter:            "Winter",
		Summer:            "Summer",
		Both:              "Both",
		W:                 "WIN",
		S:                 "SUM",
		N:                 "ER",
		ER:                "ER",
		UN:                "UNASS.",
		SearchPlaceholder: "Search...",
		SearchButton:      "Search",
		FilterButton:      "Show filters",
		ShowResults:       "Show results",
		TopFilter:         "Filter",
		ShowMore4Minus:    "Another",
		ShowMore5Plus:     "Another",
		ShowLess:          "Hide others",
		CancelFilters:     "Cancel selected filters",
		Assign:            "Assign",
		Year:              "Year",
		WinterAssign:      "Winter",
		SummerAssign:      "Summer",
		Credits:           "Credits",
		NoGuarantors:      "No guarantors",
		ReadMore:          "Read more",
		ReadLess:          "Hide",
		LoadMore:          "Show more",
		PreviousPage:      "Previous page",
		NextPage:          "Next page",
		Page:              "Page",
		Of:                "of",
		NoCoursesFound:    "No courses found.",
		// utils
		Utils: utils.Texts["en"],
	},
}
