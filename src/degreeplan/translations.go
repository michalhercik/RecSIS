package degreeplan

import (
	"github.com/michalhercik/RecSIS/utils"
	"strconv"
)

type text struct {
	Language string
	Code string
	Title string
	Status string
	Completed string
	InBlueprint string
	NotCompleted string
	Credits string
	Winter string
	Summer string
	Blueprint string
	Assign string
	Year string
	WinterAssign string
	SummerAssign string
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
		Code: "Kód",
		Title: "Název",
		Status: "Stav",
		Completed: "Splněno",
		InBlueprint: "Blueprint",
		NotCompleted: "Nesplněno",
		Credits: "Kredity",
		Winter: "ZS",
		Summer: "LS",
		Blueprint: "Blueprint",
		Assign: "Přiřadit",
		Year: "ročník",
		WinterAssign: "ZS",
		SummerAssign: "LS",
		// utils
		Utils: utils.Texts["cs"],
	},
	"en": {
		Language: "en",
		Code: "Code",
		Title: "Title",
		Status: "Status",
		Completed: "Completed",
		InBlueprint: "Blueprint",
		NotCompleted: "Not completed",
		Credits: "Credits",
		Winter: "Winter",
		Summer: "Summer",
		Blueprint: "Blueprint",
		Assign: "Assign",
		Year: "Year",
		WinterAssign: "Winter",
		SummerAssign: "Summer",
		// utils
		Utils: utils.Texts["en"],
	},
}