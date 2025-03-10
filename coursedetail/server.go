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
	router.HandleFunc(fmt.Sprintf("POST %s/like/{code}", prefix), s.like)
	router.HandleFunc(fmt.Sprintf("POST %s/dislike/{code}", prefix), s.dislike)
}

func (s Server) csPage(w http.ResponseWriter, r *http.Request) {
	s.page(w, r, texts["cs"], cs)
}

func (s Server) enPage(w http.ResponseWriter, r *http.Request) {
	s.page(w, r, texts["en"], en)
}

func (s Server) page(w http.ResponseWriter, r *http.Request, t text, lang DBLang) {
	sessionCookie, err := r.Cookie("recsis_session_key")
	if err != nil {
		log.Printf("courseDetail error: %v", err)
		return
	}
	code := r.PathValue("code")
	course, err := s.Data.Course(sessionCookie.Value, code, lang)
	if err != nil {
		log.Printf("HandlePage error %s: %v", code, err)
		PageNotFound(code, t).Render(r.Context(), w)
	} else {
		Page(course, t).Render(r.Context(), w)
	}
}

const (
	like    = 1
	dislike = 0
)

func (s Server) like(w http.ResponseWriter, r *http.Request) {
	// get the course code from the request
	sessionCookie, err := r.Cookie("recsis_session_key")
	if err != nil {
		log.Printf("like error: %v", err)
		return
	}
	code := r.PathValue("code")
	if err = s.Data.OverallRating(sessionCookie.Value, code, like); err != nil {
		log.Printf("like error: %v", err)
	}

	// TODO: should we return the updated ratings?
}

func (s Server) dislike(w http.ResponseWriter, r *http.Request) {
	// get the course code from the request
	sessionCookie, err := r.Cookie("recsis_session_key")
	if err != nil {
		log.Printf("like error: %v", err)
		return
	}
	code := r.PathValue("code")
	if err = s.Data.OverallRating(sessionCookie.Value, code, dislike); err != nil {
		log.Printf("like error: %v", err)
	}

}
