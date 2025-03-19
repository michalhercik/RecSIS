package courses

import (
	"github.com/michalhercik/RecSIS/utils"
	"strconv"
)

type text struct {
	Language string
	Winter string
	Summer string
	Both string
	W string
	S string
	N string
	ER string
	UN string
	SearchPlaceholder string
	SearchButton string
	SortByFilter string
	CreditsFilter string
	AllCredits string
	TopFilter string
	Relevance string
	Recommended string
	Rating string
	MostPopular string
	Newest string 
	SemesterFilter string
	Assign string
	Year string
	WinterAssign string
	SummerAssign string
	Credits string
	Teachers string
	ReadMore string
	ReadLess string
	LoadMore string
	PreviousPage string
	NextPage string
	Page string
	Of string
	NoCoursesFound string
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

var texts = map[string]text{
	"cs": {
		Language: "cs",
		Winter: "Zimní",
		Summer: "Letní",
		Both: "Oba",
		W: "Z",
		S: "L",
		N: "ER",
		ER: "ER",
		UN: "NE",
		SearchPlaceholder: "Hledej...",
		SearchButton: "Hledej",
		SortByFilter: "Seřadit podle",
		CreditsFilter: "Počet kreditů",
		AllCredits: "Všechny",
		TopFilter: "Filtrovat",
		Relevance: "Relevance",
		Recommended: "Doporučené",
		Rating: "Podle hodnocení",
		MostPopular: "Nejoblíbenější",
		Newest: "Nejnovější",
		SemesterFilter: "Semestr",
		Assign: "Přiřadit",
		Year: "ročník",
		WinterAssign: "ZS",
		SummerAssign: "LS",
		Credits: "Kredity",
		Teachers: "Vyučující",
		ReadMore: "Číst dále",
		ReadLess: "Sbalit",
		LoadMore: "Zobrazit další",
		PreviousPage: "Předchozí stránka",
		NextPage: "Další stránka",
		Page: "Strana",
		Of: "z",
		NoCoursesFound: "Žádné předměty nebyly nalezeny.",
		// utils
		Utils: utils.Texts["cs"],
	},
	"en": {
		Language: "en",
		Winter: "Winter",
		Summer: "Summer",
		Both: "Both",
		W: "W",
		S: "S",
		N: "N",
		ER: "ER",
		UN: "UN",
		SearchPlaceholder: "Search...",
		SearchButton: "Search",
		SortByFilter: "Sort by",
		CreditsFilter: "Credits",
		AllCredits: "All",
		TopFilter: "Filter",
		Relevance: "Relevance",
		Recommended: "Recommended",
		Rating: "By rating",
		MostPopular: "Most popular",
		Newest: "Newest",
		SemesterFilter: "Semester",
		Assign: "Assign",
		Year: "Year",
		WinterAssign: "Winter",
		SummerAssign: "Summer",
		Credits: "Credits",
		Teachers: "Teacher(s)",
		ReadMore: "Read more",
		ReadLess: "Hide",
		LoadMore: "Show more",
		PreviousPage: "Previous page",
		NextPage: "Next page",
		Page: "Page",
		Of: "of",
		NoCoursesFound: "No courses found.",
		// utils
		Utils: utils.Texts["en"],
	},
}