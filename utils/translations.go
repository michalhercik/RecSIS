package utils

type Text struct {
	Language string
	Home string
	Courses string
	Blueprint string
	DegreePlan string 
	Login string
	Contact string
}

func (t Text) LangLink(URL string) string {
	return "/" + t.Language + URL
}

var Texts = map[string]Text{
	"cs": {
		Language: "cs",
		Home: "Domů",
		Courses: "Předměty",
		Blueprint: "Blueprint",
		DegreePlan: "Studijní plán",
		Login: "Přihlášení",
		Contact: "V případě jakýchkoliv problémů kontaktujte tým RecSIS.",
	},
	"en": {
		Language: "en",
		Home: "Home",
		Courses: "Courses",
		Blueprint: "Blueprint",
		DegreePlan: "Degree plan",
		Login: "Login",
		Contact: "In case of any problems, please contact the RecSIS team.",
	},
}