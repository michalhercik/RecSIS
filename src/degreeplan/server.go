package degreeplan

import (
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/language"
)

//================================================================================
// Server Type
//================================================================================

type Server struct {
	router   *http.ServeMux
	Auth     Authentication
	BpBtn    BlueprintAddButton
	Data     DBManager
	DPSearch MeiliSearch
	Page     Page
}

func (s Server) Router() http.Handler {
	return s.router
}

type Authentication interface {
	UserID(r *http.Request) (string, error)
}

type BlueprintAddButton interface {
	PartialComponent(lang language.Language) PartialBlueprintAdd
	PartialComponentSecond(lang language.Language) PartialBlueprintAdd
	ParseRequest(r *http.Request) ([]string, int, int, error)
	Action(userID string, year int, semester int, course ...string) ([]int, error)
}

type PartialBlueprintAdd = func(hxSwap, hxTarget, hxInclude string, years []bool, course ...string) templ.Component

type Page interface {
	View(main templ.Component, lang language.Language, title string, userID string) templ.Component
}

//================================================================================
// Routing
//================================================================================

func (s *Server) Init() {
	router := http.NewServeMux()
	router.HandleFunc("GET /", s.page)
	router.HandleFunc("GET /{dpCode}", s.show)
	router.HandleFunc("GET /search", s.searchDegreePlan)
	router.HandleFunc("POST /blueprint", s.addCourseToBlueprint)
	s.router = router
}

//================================================================================
// Handlers
//================================================================================

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	s.renderPage(w, r, userID, lang)
}

func (s Server) show(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	t := texts[lang]
	dpCode := r.PathValue("dpCode")
	dpYear, err := strconv.Atoi(r.FormValue("dp-year"))
	if err != nil {
		http.Error(w, "Invalid degree plan year", http.StatusBadRequest)
		log.Printf("renderPage: %v", err)
		return
	}
	dp, err := s.Data.degreePlan(userID, dpCode, dpYear, lang)
	if err != nil {
		http.Error(w, "Unable to retrieve degree plan", http.StatusInternalServerError)
		log.Printf("renderPage: %v", err)
		return
	}
	partialBpBtn := s.BpBtn.PartialComponent(lang)
	partialBpBtnChecked := s.BpBtn.PartialComponentSecond(lang)
	main := Content(dp, t, partialBpBtn, partialBpBtnChecked)
	s.Page.View(main, lang, t.pageTitle, userID).Render(r.Context(), w)
}

func (s Server) searchDegreePlan(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	query := r.FormValue("q")
	results, err := s.DPSearch.QuickSearch(quickRequest{
		query: query,
		limit: 5,
	})
	if err != nil {
		http.Error(w, "Search failed", http.StatusInternalServerError)
		log.Printf("searchDegreePlan error: %v", err)
		return
	}
	QuickSearchResultsContent(results.DegreePlans, texts[lang]).Render(r.Context(), w)
}

func (s Server) addCourseToBlueprint(w http.ResponseWriter, r *http.Request) {
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	courseCode, year, semester, err := s.BpBtn.ParseRequest(r)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		log.Printf("HandlePage error: %v", err)
		return
	}
	courseCode = append(courseCode, r.Form["selected"]...)
	_, err = s.BpBtn.Action(userID, year, semester, courseCode...)
	if err != nil {
		http.Error(w, "Unable to add course to blueprint", http.StatusInternalServerError)
		log.Printf("HandlePage error: %v", err)
		return
	}
	s.renderContent(w, r, userID, language.FromContext(r.Context()))
}

func (s Server) renderPage(w http.ResponseWriter, r *http.Request, userID string, lang language.Language) {
	t := texts[lang]
	dp, err := s.Data.userDegreePlan(userID, lang)
	if err != nil {
		http.Error(w, "Unable to retrieve degree plan", http.StatusInternalServerError)
		log.Printf("renderPage: %v", err)
		return
	}
	partialBpBtn := s.BpBtn.PartialComponent(lang)
	partialBpBtnChecked := s.BpBtn.PartialComponentSecond(lang)
	main := Content(dp, t, partialBpBtn, partialBpBtnChecked)
	s.Page.View(main, lang, t.pageTitle, userID).Render(r.Context(), w)
}

func (s Server) renderContent(w http.ResponseWriter, r *http.Request, userID string, lang language.Language) {
	t := texts[lang]
	dp, err := s.Data.userDegreePlan(userID, lang)
	if err != nil {
		http.Error(w, "Unable to retrieve degree plan", http.StatusInternalServerError)
		log.Printf("renderPage: %v", err)
		return
	}
	partialBpBtn := s.BpBtn.PartialComponent(lang)
	partialBpBtnChecked := s.BpBtn.PartialComponentSecond(lang)
	Content(dp, t, partialBpBtn, partialBpBtnChecked).Render(r.Context(), w)
}
