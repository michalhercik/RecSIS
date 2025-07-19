package filters

import (
	"github.com/michalhercik/RecSIS/language"
)

type text struct {
	errCategoryNotFound string
	errValueNotFound    string
}

var texts = map[language.Language]text{
	language.CS: {
		errCategoryNotFound: "kategorie nenalezena",
		errValueNotFound:    "hodnota nenalezena v kategorii",
	},
	language.EN: {
		errCategoryNotFound: "Category not found",
		errValueNotFound:    "Value not found in category",
	},
}
