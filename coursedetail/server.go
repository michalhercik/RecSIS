package coursedetail

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	Data DataManager
}

func (s Server) Register(router *http.ServeMux, prefix string) {
	//router.HandleFunc(fmt.Sprintf("GET %s/{code}", prefix), s.page) // TODO get language from http header
	router.HandleFunc(fmt.Sprintf("GET /cs%s/{code}", prefix), s.csPage)
	router.HandleFunc(fmt.Sprintf("GET /en%s/{code}", prefix), s.enPage)
	// TODO: should we differentiate between languages for POSTs?
	router.HandleFunc(fmt.Sprintf("POST %s/{code}/like", prefix), s.like)
	router.HandleFunc(fmt.Sprintf("POST %s/{code}/dislike", prefix), s.dislike)
}

func (s Server) csPage(w http.ResponseWriter, r *http.Request) {
	s.page(w, r, texts["cs"], cs)
}

func (s Server) enPage(w http.ResponseWriter, r *http.Request) {
	s.page(w, r, texts["en"], en)
}

func (s Server) page(w http.ResponseWriter, r *http.Request, t text, lang DBLang) {
	code := r.PathValue("code")
	course, err := s.Data.Course(code, lang)
	if err != nil {
		log.Printf("HandlePage error %s: %v", code, err)
		PageNotFound(code, t).Render(r.Context(), w)
	} else {
		Page(course, t).Render(r.Context(), w)
	}
}

func (s Server) like(w http.ResponseWriter, r *http.Request) {
	// get the course code from the request
	code := r.PathValue("code")

	// TODO implement the like functionality

	Ratings([]Rating{}, code).Render(r.Context(), w)
}

func (s Server) dislike(w http.ResponseWriter, r *http.Request) {
	// get the course code from the request
	code := r.PathValue("code")

	// TODO implement the dislike functionality

	Ratings([]Rating{}, code).Render(r.Context(), w)
}
