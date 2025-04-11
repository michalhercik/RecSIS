package home

import (
	"log"
	"net/http"

	"github.com/michalhercik/RecSIS/language"
)

type Server struct {
	router *http.ServeMux
}

func (s *Server) Init() {
	router := http.NewServeMux()
	router.HandleFunc("GET /", s.page)
	router.HandleFunc("GET /home/", s.page)
	s.router = router
}

func (s Server) Router() http.Handler {
	return s.router
}

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	log.Println("Home page")
	log.Println(r.URL.Path)
	lang := language.FromContext(r.Context())
	t := texts[lang]
	Page(t).Render(r.Context(), w)
}
