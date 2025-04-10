package home

import (
	"net/http"

	"github.com/michalhercik/RecSIS/language"
)

type Server struct{}

func (s Server) Register(router *http.ServeMux) {
	lr := language.LanguageRouter{Router: router}
	lr.HandleLangFunc("/home/", "GET", s.page)
	lr.HandleLangFunc("/", "GET", s.page)
}

func (s Server) page(w http.ResponseWriter, r *http.Request, lang language.Language) {
	t := texts[lang]
	Page(t).Render(r.Context(), w)
}
