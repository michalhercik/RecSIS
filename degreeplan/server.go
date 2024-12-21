package degreeplan

import (
	"fmt"
	"net/http"
)

type Server struct{}

func (s Server) Register(router *http.ServeMux, prefix string) {
	router.HandleFunc(fmt.Sprintf("GET %s", prefix), s.page)
}

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	Page().Render(r.Context(), w)
}
