package degreeplan

import (
	"log"
	"net/http"

	"github.com/michalhercik/RecSIS/language"
)

type Server struct {
	router *http.ServeMux
	Data   DataManager
}

func (s *Server) Init() {
	router := http.NewServeMux()
	router.HandleFunc("GET /", s.page)
	s.router = router
}

func (s Server) Router() http.Handler {
	return s.router
}

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
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
