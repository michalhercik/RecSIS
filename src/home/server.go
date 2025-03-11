package home

import (
	"net/http"
)

type Server struct{}

func (s Server) Register(router *http.ServeMux) {
	//router.HandleFunc("GET /{$}", s.page) // TODO get language from http header
	router.HandleFunc("GET /cs", s.csPage)
	router.HandleFunc("GET /en", s.enPage)
	//router.HandleFunc("GET /home", s.page) // TODO get language from http header
	router.HandleFunc("GET /cs/home", s.csPage)
	router.HandleFunc("GET /en/home", s.enPage)
}

// TODO struct to json

func (s Server) csPage(w http.ResponseWriter, r *http.Request) {
	s.page(w, r, texts["cs"])
}

func (s Server) enPage(w http.ResponseWriter, r *http.Request) {
	s.page(w, r, texts["en"]) 
}

func (s Server) page(w http.ResponseWriter, r *http.Request, t text) {
	Page(t).Render(r.Context(), w)
}


