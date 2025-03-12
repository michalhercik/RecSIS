package degreeplan

import (
	"github.com/michalhercik/RecSIS/utils"
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
	// utils
	Utils utils.Text
}

var texts = map[string]text{
	"cs": {
		Language: "cs",
		Code: "Kód",
		Title: "Název",
		Status: "Stav",
		Completed: "Splněno",
		InBlueprint: "V blueprintu",
		NotCompleted: "Nesplněno",
		Credits: "Kredity",
		Winter: "ZS",
		Summer: "LS",
		Blueprint: "Blueprint",
		// utils
		Utils: utils.Texts["cs"],
	},
	"en": {
		Language: "en",
		Code: "Code",
		Title: "Title",
		Status: "Status",
		Completed: "Completed",
		InBlueprint: "In blueprint",
		NotCompleted: "Not completed",
		Credits: "Credits",
		Winter: "Winter",
		Summer: "Summer",
		Blueprint: "Blueprint",
		// utils
		Utils: utils.Texts["en"],
	},
}