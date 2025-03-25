package coursedetail

import (
	"github.com/michalhercik/RecSIS/utils"
	"strconv"
)

type text struct {
	Language string
	Faculty string
	Semester string
	Winter string
	Summer string
	Both string
	WinterAssign string
	SummerAssign string
	BothAssign string
	Year string
	Assign string
	ECredits string
	Capacity string
	CapacityNoLimit string
	ScopeExam string
	StateOfCourse string
	LanguageOfCourse string
	AdditionalInfo string
	Guarantors string
	Teachers string
	SISLink string
	Comments string
	CourseWithCode string
	NotFound string
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
		Faculty: "Fakulta",
		Semester: "Semestr",
		Winter: "ZS",
		Summer: "LS",
		Both: "oba",
		WinterAssign: "ZS",
		SummerAssign: "LS",
		BothAssign: "Oba",
		Year: "ročník",
		Assign: "Přiřadit",
		ECredits: "E-Kredity",
		Capacity: "Kapacita",
		CapacityNoLimit: "bez omezení",
		ScopeExam: "Rozsah, examinace",
		StateOfCourse: "Stav předmětu",
		LanguageOfCourse: "Jazyk výuky",
		AdditionalInfo: "Další informace",
		Guarantors: "Garant(i)",
		Teachers: "Vyučující",
		SISLink: "Odkaz do SIS",
		Comments: "Komentáře",
		CourseWithCode: "Předmět s kódem ",
		NotFound: " nenalezen.",
		// utils
		Utils: utils.Texts["cs"],
	},
	"en": {
		Language: "en",
		Faculty: "Faculty",
		Semester: "Semester",
		Winter: "Winter",
		Summer: "Summer",
		Both: "Both",
		WinterAssign: "Winter",
		SummerAssign: "Summer",
		BothAssign: "Both",
		Year: "Year",
		Assign: "Assign",
		ECredits: "E-Credits",
		Capacity: "Capacity",
		CapacityNoLimit: "No limit",
		ScopeExam: "Scope, examination",
		StateOfCourse: "State of the course",
		LanguageOfCourse: "Language",
		AdditionalInfo: "Additional information",
		Guarantors: "Guarantor(s)",
		Teachers: "Teacher(s)",
		SISLink: "Link to SIS",
		Comments: "Comments",
		CourseWithCode: "Course with code ",
		NotFound: " not found.",
		// utils
		Utils: utils.Texts["en"],
	},
}