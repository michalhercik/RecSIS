package degreeplan

import (
	"log"
	"net/http"

	"github.com/michalhercik/RecSIS/language"
)

type Server struct {
	router *http.ServeMux
	Data   DBManager
	Auth   Authentication
	BpBtn  BlueprintAddButton
	Page   Page
}

func (s *Server) Init() {
	router := http.NewServeMux()
	router.HandleFunc("GET /", s.page)
	router.HandleFunc("POST /blueprint", s.addCourseToBlueprint)
	s.router = router
}

func (s Server) Router() http.Handler {
	return s.router
}

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	s.renderPage(w, r, userID, lang)
}

func (s Server) addCourseToBlueprint(w http.ResponseWriter, r *http.Request) {
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// TODO: should we differentiate between BpBtnRow and BpBtnChecked here?
	courseCode, year, semester, err := s.BpBtn.ParseRequest(r)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		log.Printf("HandlePage error: %v", err)
		return
	}
	courseCode = append(courseCode, r.Form["selected"]...)
	// courseCode := r.Form["selected"]
	// TODO: should we differentiate between BpBtn and BpBtnChecked here?
	_, err = s.BpBtn.Action(userID, year, semester, courseCode...)
	if err != nil {
		http.Error(w, "Unable to add course to blueprint", http.StatusInternalServerError)
		log.Printf("HandlePage error: %v", err)
		return
	}
	s.renderPage(w, r, userID, language.FromContext(r.Context()))
}

func (s Server) renderPage(w http.ResponseWriter, r *http.Request, userID string, lang language.Language) {
	t := texts[lang]
	dp, err := s.Data.DegreePlan(userID, lang)
	if err != nil {
		http.Error(w, "Unable to retrieve degree plan", http.StatusInternalServerError)
		log.Printf("renderPage: %v", err)
		return
	}
	partialBpBtn := s.BpBtn.PartialComponent(lang)
	partialBpBtnChecked := s.BpBtn.PartialComponentSecond(lang)
	main := Content(dp, t, partialBpBtn, partialBpBtnChecked)
	s.Page.View(main, lang, t.PageTitle).Render(r.Context(), w)
}
