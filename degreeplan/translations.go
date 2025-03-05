package degreeplan

import (
	"github.com/michalhercik/RecSIS/utils"
)

type text struct {
	Language string
	Code string
	Title string
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
		Credits: "Credits",
		Winter: "Winter",
		Summer: "Summer",
		Blueprint: "Blueprint",
		// utils
		Utils: utils.Texts["en"],
	},
}