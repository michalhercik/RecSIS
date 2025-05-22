package degreeplan

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/michalhercik/RecSIS/dbds"
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
	router.HandleFunc("POST /blueprint/{coursecode}", s.addCourseToBlueprint)
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

	time.Sleep(2 * time.Second) // Simulate a delay for testing

	s.renderPage(w, r, userID, lang)
}

func (s Server) addCourseToBlueprint(w http.ResponseWriter, r *http.Request) {
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	courseCode := r.PathValue("coursecode")
	year, err := strconv.Atoi(r.FormValue("year"))
	if err != nil {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	semesterInt, err := strconv.Atoi(r.FormValue("semester"))
	if err != nil {
		http.Error(w, "Invalid semester", http.StatusBadRequest)
		return
	}
	semester := dbds.SemesterAssignment(semesterInt)
	_, err = s.BpBtn.Action(userID, courseCode, year, semester)
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
	numberOfYears, err := s.BpBtn.NumberOfYears(userID)
	if err != nil {
		http.Error(w, "Unable to retrieve number of years", http.StatusInternalServerError)
		log.Printf("renderPage: %v", err)
		return
	}
	partialComponent := s.BpBtn.PartialComponent(numberOfYears, lang)
	main := Content(dp, t, partialComponent)
	s.Page.View(main, lang, t.PageTitle).Render(r.Context(), w)
}
