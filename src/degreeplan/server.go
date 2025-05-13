package degreeplan

import (
	"log"
	"net/http"
	"strconv"

	"github.com/michalhercik/RecSIS/dbcourse"
	"github.com/michalhercik/RecSIS/language"
)

type Server struct {
	router *http.ServeMux
	Data   DataManager
	Auth   Authentication
	BpBtn  BlueprintAddButton
	Page   Page
}

func (s *Server) Init() {
	router := http.NewServeMux()
	router.HandleFunc("GET /", s.page)
	router.HandleFunc("POST /blueprint/{coursecode}", s.AddCourseToBlueprint)
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
	numberOfYears, err := s.BpBtn.NumberOfYears(userID)
	if err != nil {
		http.Error(w, "Unable to retrieve number of years", http.StatusInternalServerError)
		log.Printf("HandlePage error: %v", err)
		return
	}
	partialComponent := s.BpBtn.PartialComponent(numberOfYears, lang)
	main := Content(dp, t, partialComponent)
	s.Page.View(main, lang, t.PageTitle).Render(r.Context(), w)
}

func (s Server) AddCourseToBlueprint(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
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
	semester := dbcourse.SemesterAssignment(semesterInt)
	_, err = s.BpBtn.Action(userID, courseCode, year, semester)
	if err != nil {
		http.Error(w, "Unable to add course to blueprint", http.StatusInternalServerError)
		log.Printf("HandlePage error: %v", err)
		return
	}
	numberOfYears, err := s.BpBtn.NumberOfYears(userID)
	if err != nil {
		http.Error(w, "Unable to retrieve number of years", http.StatusInternalServerError)
		log.Printf("HandlePage error: %v", err)
		return
	}
	t := texts[lang]
	btn := s.BpBtn.PartialComponent(numberOfYears, lang)
	course, err := s.Data.Course(userID, courseCode, lang)
	if err != nil {
		http.Error(w, "Unable to retrieve course", http.StatusInternalServerError)
		log.Printf("HandlePage error: %v", err)
		return
	}
	CourseRow(&course, btn, t).Render(r.Context(), w)
}
