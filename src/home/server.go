package home

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/language"
)

//================================================================================
// Server Type
//================================================================================

type Server struct {
	Auth        Authentication
	Page        Page
	Recommender string
	router      *http.ServeMux
}

func (s Server) Router() http.Handler {
	return s.router
}

type Authentication interface {
	UserID(r *http.Request) (string, error)
}

type Page interface {
	View(main templ.Component, lang language.Language, title string, userID string) templ.Component
}

//================================================================================
// Routing
//================================================================================

func (s *Server) Init() {
	router := http.NewServeMux()
	router.HandleFunc("GET /", s.page)
	router.HandleFunc("GET /home/", s.page)
	s.router = router
}

//================================================================================
// Handlers
//================================================================================

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]

	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	recommended, _ := s.fetchCourses("recommended", lang)
	newest, _ := s.fetchCourses("newest", lang)

	content := homePage{
		recommendedCourses: recommended,
		newCourses:         newest,
	}

	main := Content(&content, t)
	s.Page.View(main, lang, t.pageTitle, userID).Render(r.Context(), w)
}

func (s Server) fetchCourses(endpoint string, lang language.Language) ([]course, error) {
	url := fmt.Sprintf("%s/%s?lang=%s", s.Recommender, endpoint, lang)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var courses []course
	err = json.Unmarshal(body, &courses)
	return courses, err
}
