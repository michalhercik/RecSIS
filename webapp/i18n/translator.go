package i18n

import "github.com/michalhercik/RecSIS/language"

type Translatable string

const (
	ViewAll Translatable = "viewAll"
	Credits Translatable = "credits"
)

var translations = map[Translatable]language.LangString{
	ViewAll: {
		EN: "View all",
		CS: "Zobrazit vše",
	},
	Credits: {
		EN: "Cr.",
		CS: "Kr.",
	},
}

type Translator struct {
	lang language.Language
}

func New(lang language.Language) Translator {
	return Translator{lang: lang}
}

func (t Translator) T(key Translatable) string {
	return translations[key].String(t.lang)
}

func TranslatorFunc(lang language.Language) func(key Translatable) string {
	translator := New(lang)
	return func(key Translatable) string {
		return translator.T(key)
	}
}
