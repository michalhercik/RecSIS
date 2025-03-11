package blueprint

import (
	"github.com/michalhercik/RecSIS/utils"
	"strconv"
)

type text struct {
	Language string
	NumOfYears string
	Total string
	Code string
	Title string
	Credits string
	Winter string
	WinterLong string
	Summer string
	SummerLong string
	Teachers string
	Unassigned string
	Year string
	YearBig string
	Semester string
	// tooltips
	TTUnassignChecked	string
	TTAssignChecked string
	TTRemoveChecked string
	TTRemoveUnassigned string
	TTAssign string
	TTRemove string
	TTMove string
	TTUnassignYear1 string
	TTUnassignYear2 string
	TTRemoveYear1 string
	TTRemoveYear2 string
	TTUnassignWinter string
	TTRemoveWinter string
	TTUnassign string
	TTReassign string
	TTUnassignSummer string
	TTRemoveSummer string
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
		NumOfYears: "Počet ročníků",
		Total: "Celkem",
		Code: "Kód",
		Title: "Název",
		Credits: "Kredity",
		Winter: "ZS",
		WinterLong: "Zimní",
		Summer: "LS",
		SummerLong: "Letní",
		Teachers: "Vyučující",
		Unassigned: "Nezařazené",
		Year: "ročník",
		YearBig: "Ročník",
		Semester: "Semestr",
		// tooltips
		TTUnassignChecked: "Přesunout vybrané předměty do nezařazených",
		TTAssignChecked: "Zařadit vybrané předměty",
		TTRemoveChecked: "Odstranit vybrané předměty",
		TTRemoveUnassigned: "Odstranit nezařazené předměty",
		TTAssign: "Zařadit předmět",
		TTRemove: "Odstranit předmět",
		TTMove: "Přesunout předmět pomocí drag-and-drop",
		TTUnassignYear1: "Přesunout všechny předměty z ",
		TTUnassignYear2: ". ročníku do nezařazených",
		TTRemoveYear1: "Odstranit všechny předměty z ",
		TTRemoveYear2: ". ročníku",
		TTUnassignWinter: "Přesunout všechny předměty z tohoto zimního semestru do nezařazených",
		TTRemoveWinter: "Odstranit všechny předměty z tohoto zimního semestru",
		TTUnassign: "Přesunout předmět do nezařazených",
		TTReassign: "Přesunout předmět",
		TTUnassignSummer: "Přesunout všechny předměty z tohoto letního semestru do nezařazených",
		TTRemoveSummer: "Odstranit všechny předměty z tohoto letního semestru",
		// utils
		Utils: utils.Texts["cs"],
	},
	"en": {
		Language: "en",
		NumOfYears: "Number of years",
		Total: "Total",
		Code: "Code",
		Title: "Title",
		Credits: "Credits",
		Winter: "Winter",
		WinterLong: "Winter",
		Summer: "Summer",
		SummerLong: "Summer",
		Teachers: "Teacher(s)",
		Unassigned: "Unassigned",
		Year: "Year",
		YearBig: "Year",
		Semester: "Semester",
		// tooltips
		TTUnassignChecked: "Unassign all selected courses",
		TTAssignChecked: "Assign all selected courses",
		TTRemoveChecked: "Remove all selected courses",
		TTRemoveUnassigned: "Remove all unassigned courses",
		TTAssign: "Assign course",
		TTRemove: "Remove course",
		TTMove: "Drag and drop to sort",
		TTUnassignYear1: "Unassign all courses from Year ",
		TTUnassignYear2: "",
		TTRemoveYear1: "Remove all courses in Year ",
		TTRemoveYear2: "",
		TTUnassignWinter: "Unassign all courses from this winter semester",
		TTRemoveWinter: "Remove all courses from this winter semester",
		TTUnassign: "Unassign course",
		TTReassign: "Reassign course",
		TTUnassignSummer: "Unassign all courses from this summer semester",
		TTRemoveSummer: "Remove all courses from this summer semester",
		// utils
		Utils: utils.Texts["en"],
	},
}