package home

import (
	"net/http"
)

type Server struct{}

func (s Server) Register(router *http.ServeMux) {
	router.HandleFunc("GET /{$}", s.page)
	router.HandleFunc("GET /home", s.page)
}

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	Page().Render(r.Context(), w)
}
