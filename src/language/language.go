package language

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Language string

const (
	CS Language = "cs"
	EN Language = "en"
)

type LanguageFuncHandler func(http.ResponseWriter, *http.Request, Language)

func DefaultLanguage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !(strings.HasPrefix(r.URL.Path, "/"+string(CS)) || strings.HasPrefix(r.URL.Path, "/"+string(EN))) {
			path, err := url.JoinPath(string(CS), r.URL.Path)
			if err != nil {
				log.Println("DefaultLanguage:", err)
				return
			}
			r.URL.Path = path
		}
		next.ServeHTTP(w, r)
	})
}

type LanguageRouter struct {
	Router *http.ServeMux
}

func (lr LanguageRouter) HandleLangFunc(path string, method string, handler LanguageFuncHandler) {
	lr.Router.HandleFunc(method+" /"+string(CS)+path, CS.Handle(handler))
	lr.Router.HandleFunc(method+" /"+string(EN)+path, EN.Handle(handler))
}

func (l Language) Handle(handler LanguageFuncHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, l)
	}
}
