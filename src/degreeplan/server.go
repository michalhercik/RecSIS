package degreeplan

import (
	"log"
	"net/http"

	"github.com/michalhercik/RecSIS/language"
)

type Server struct {
	router *http.ServeMux
	Data   DataManager
	Auth   Authentication
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
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	dp, err := s.Data.DegreePlan(userID, lang)
	if err != nil {
		http.Error(w, "Unable to retrieve degree plan", http.StatusInternalServerError)
		log.Printf("HandlePage error: %v", err)
		return
	}
	Page(dp, t).Render(r.Context(), w)
}
