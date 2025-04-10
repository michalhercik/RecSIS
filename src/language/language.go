package language

import "net/http"

type Language string

const (
	CS Language = "cs"
	EN Language = "en"
)

type LanguageFuncHandler func(http.ResponseWriter, *http.Request, Language)

func (l Language) Handle(handler LanguageFuncHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, l)
	}
}

type LanguageRouter struct {
	Router *http.ServeMux
}

func (lr LanguageRouter) HandleLangFunc(path string, method string, handler LanguageFuncHandler) {
	lr.Router.HandleFunc(method+" /cs"+path, CS.Handle(handler))
	lr.Router.HandleFunc(method+" /en"+path, EN.Handle(handler))
}

func (lr LanguageRouter) HandleFunc(path string, method string, handler func(http.ResponseWriter, *http.Request)) {
	lr.Router.HandleFunc(method+" "+path, handler)
}
