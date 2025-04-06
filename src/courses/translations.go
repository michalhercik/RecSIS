package courses

import (
	"strconv"

	"github.com/michalhercik/RecSIS/utils"
)

type text struct {
	Language               string
	Winter                 string
	Summer                 string
	Both                   string
	W                      string
	S                      string
	N                      string
	ER                     string
	UN                     string
	SearchPlaceholder      string
	SearchButton           string
	CreditsFilter          string
	FacultyFilter          string
	SemesterFilter         string
	ExamTypeFilter         string
	SemesterCountFilter    string
	LectureRangeWFilter    string
	LectureRangeSFilter    string
	SeminarRangeWFilter    string
	SeminarRangeSFilter    string
	RangeUnitFilter        string
	TaughtFilter           string
	LanguageFilter         string
	FacultyGuarantorFilter string
	CapacityFilter         string
	MinNumberFilter        string
	TopFilter              string
	Assign                 string
	Year                   string
	WinterAssign           string
	SummerAssign           string
	Credits                string
	Teachers               string
	ReadMore               string
	ReadLess               string
	LoadMore               string
	PreviousPage           string
	NextPage               string
	Page                   string
	Of                     string
	NoCoursesFound         string
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
		Language:               "cs",
		Winter:                 "Zimní",
		Summer:                 "Letní",
		Both:                   "Oba",
		W:                      "Z",
		S:                      "L",
		N:                      "ER",
		ER:                     "ER",
		UN:                     "NE",
		SearchPlaceholder:      "Hledej...",
		SearchButton:           "Hledej",
		CreditsFilter:          "Počet kreditů",
		FacultyFilter:          "Fakulta",
		SemesterFilter:         "Semestr",
		ExamTypeFilter:         "Typ zkoušky",
		SemesterCountFilter:    "Počet semestrů",
		LectureRangeWFilter:    "Rozsah přednášek (ZS)",
		LectureRangeSFilter:    "Rozsah přednášek (LS)",
		SeminarRangeWFilter:    "Rozsah cvičení (ZS)",
		SeminarRangeSFilter:    "Rozsah cvičení (LS)",
		RangeUnitFilter:        "Jednotka rozsahu",
		TaughtFilter:           "Vyučováno",
		LanguageFilter:         "Jazyk(y)",
		FacultyGuarantorFilter: "Garant fakulty???TODO",
		CapacityFilter:         "Kapacita",
		MinNumberFilter:        "Minimální počet zapsaných",
		TopFilter:              "Filtrovat",
		Assign:                 "Přiřadit",
		Year:                   "ročník",
		WinterAssign:           "ZS",
		SummerAssign:           "LS",
		Credits:                "Kredity",
		Teachers:               "Vyučující",
		ReadMore:               "Číst dále",
		ReadLess:               "Sbalit",
		LoadMore:               "Zobrazit další",
		PreviousPage:           "Předchozí stránka",
		NextPage:               "Další stránka",
		Page:                   "Strana",
		Of:                     "z",
		NoCoursesFound:         "Žádné předměty nebyly nalezeny.",
		// utils
		Utils: utils.Texts["cs"],
	},
	"en": {
		Language:               "en",
		Winter:                 "Winter",
		Summer:                 "Summer",
		Both:                   "Both",
		W:                      "W",
		S:                      "S",
		N:                      "N",
		ER:                     "ER",
		UN:                     "UN",
		SearchPlaceholder:      "Search...",
		SearchButton:           "Search",
		CreditsFilter:          "Credits",
		FacultyFilter:          "Faculty",
		SemesterFilter:         "Semester",
		ExamTypeFilter:         "Exam type",
		SemesterCountFilter:    "Semester count",
		LectureRangeWFilter:    "Lecture range (Winter)",
		LectureRangeSFilter:    "Lecture range (Summer)",
		SeminarRangeWFilter:    "Seminar range (Winter)",
		SeminarRangeSFilter:    "Seminar range (Summer)",
		RangeUnitFilter:        "Unit of range",
		TaughtFilter:           "Taught",
		LanguageFilter:         "Language(s)",
		FacultyGuarantorFilter: "Faculty guarantor???TODO",
		CapacityFilter:         "Capacity",
		MinNumberFilter:        "Minimum number of students",
		TopFilter:              "Filter",
		Assign:                 "Assign",
		Year:                   "Year",
		WinterAssign:           "Winter",
		SummerAssign:           "Summer",
		Credits:                "Credits",
		Teachers:               "Teacher(s)",
		ReadMore:               "Read more",
		ReadLess:               "Hide",
		LoadMore:               "Show more",
		PreviousPage:           "Previous page",
		NextPage:               "Next page",
		Page:                   "Page",
		Of:                     "of",
		NoCoursesFound:         "No courses found.",
		// utils
		Utils: utils.Texts["en"],
	},
}
