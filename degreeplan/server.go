package degreeplan

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	Data DataManager
}

func (s Server) Register(router *http.ServeMux, prefix string) {
	//router.HandleFunc(fmt.Sprintf("GET /%s", prefix), s.page) // TODO get language from http header
	router.HandleFunc(fmt.Sprintf("GET /cs%s", prefix), s.csPage)
	router.HandleFunc(fmt.Sprintf("GET /en%s", prefix), s.enPage)
}

func (s Server) csPage(w http.ResponseWriter, r *http.Request) {
	s.page(w, r, cs, texts["cs"])
}

func (s Server) enPage(w http.ResponseWriter, r *http.Request) {
	s.page(w, r, en, texts["en"])
}

func (s Server) page(w http.ResponseWriter, r *http.Request, lang DBLang, t text) {
	// TODO: redirect to login page?
	uid_cookie, err := r.Cookie("recsis_session_key")
	if err != nil {
		http.Error(w, "unknown student", http.StatusBadRequest)
		log.Printf("HandlePage error: %v", err)
		return
	}
	dp, err := s.Data.DegreePlan(uid_cookie.Value, lang)
	if err != nil {
		http.Error(w, "Unable to retrieve degree plan", http.StatusInternalServerError)
		log.Printf("HandlePage error: %v", err)
		return
	}
	Page(dp, t).Render(r.Context(), w)
}
