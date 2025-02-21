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
	router.HandleFunc(fmt.Sprintf("GET %s", prefix), s.page)
}

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	// TODO: redirect to login page?
	uid_cookie, err := r.Cookie("recsis_session_key")
	if err != nil {
		http.Error(w, "unknown student", http.StatusBadRequest)
		log.Printf("HandlePage error: %v", err)
		return
	}
	dp, err := s.Data.DegreePlan(uid_cookie.Value, cs)
	if err != nil {
		http.Error(w, "Unable to retrieve degree plan", http.StatusInternalServerError)
		log.Printf("HandlePage error: %v", err)
		return
	}
	Page(dp).Render(r.Context(), w)
}
