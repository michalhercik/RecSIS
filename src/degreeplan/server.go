package degreeplan

import (
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/language"
)

//================================================================================
// Server Type
//================================================================================

type Server struct {
	Auth     Authentication
	BpBtn    BlueprintAddButton
	Data     DBManager
	DPSearch MeiliSearch
	Error    Error
	Page     Page
	router   http.Handler
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
	ParseRequest(r *http.Request, additionalCourses []string) ([]string, int, int, error)
	Action(userID string, year int, semester int, lang language.Language, course ...string) ([]int, error)
}

type PartialBlueprintAdd = func(hxSwap, hxTarget, hxInclude string, years []bool, course string) templ.Component

type Error interface {
	Log(err error)
	Render(w http.ResponseWriter, r *http.Request, code int, userMsg string, lang language.Language)
	RenderPage(w http.ResponseWriter, r *http.Request, code int, userMsg string, title string, userID string, lang language.Language)
}

type Page interface {
	View(main templ.Component, lang language.Language, title string, userID string) templ.Component
}

//================================================================================
// Routing
//================================================================================

func (s *Server) Init() {
	router := http.NewServeMux()
	router.HandleFunc("GET /{$}", s.page)
	router.HandleFunc("GET /show/{dpCode}", s.show)
	router.HandleFunc("GET /search", s.searchDegreePlan)
	router.HandleFunc("POST /blueprint", s.addCourseToBlueprint)

	// Wrap mux to catch unmatched routes
	s.router = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if mux has a handler for the URL
		_, pattern := router.Handler(r)
		if pattern == "" {
			s.pageNotFound(w, r)
			return
		}
		router.ServeHTTP(w, r)
	})
}

//================================================================================
// Handlers
//================================================================================

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID, err := s.Auth.UserID(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, "", lang)
		return
	}
	s.renderPage(w, r, userID, lang)
}

func (s Server) show(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID, err := s.Auth.UserID(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, "", lang)
		return
	}
	dpYearString := r.FormValue(searchDegreePlanYear)
	dpYear, err := strconv.Atoi(dpYearString)
	if err != nil {
		s.Error.Log(errorx.AddContext(err, errorx.P(searchDegreePlanYear, dpYearString)))
		s.Error.RenderPage(w, r, http.StatusBadRequest, t.errInvalidDPYear, t.pageTitle, userID, lang)
		return
	}
	dpCode := r.PathValue("dpCode")
	dp, err := s.Data.degreePlan(userID, dpCode, dpYear, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, userID, lang)
		return
	}
	partialBpBtn := s.BpBtn.PartialComponent(lang)
	partialBpBtnChecked := s.BpBtn.PartialComponentSecond(lang)
	main := Content(dp, t, partialBpBtn, partialBpBtnChecked)
	s.Page.View(main, lang, t.pageTitle, userID).Render(r.Context(), w)
}

func (s Server) searchDegreePlan(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	query := r.FormValue(searchDegreePlanName)
	results, err := s.DPSearch.QuickSearch(quickRequest{
		query: query,
		limit: searchDegreePlanLimit,
	}, t)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	QuickSearchResultsContent(results.DegreePlans, t).Render(r.Context(), w)
}

func (s Server) addCourseToBlueprint(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID, err := s.Auth.UserID(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, "", lang)
		return
	}
	r.ParseForm()
	additionalCourses := r.Form[checkboxName]
	courseCode, year, semester, err := s.BpBtn.ParseRequest(r, additionalCourses)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	_, err = s.BpBtn.Action(userID, year, semester, lang, courseCode...)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	s.renderContent(w, r, userID, language.FromContext(r.Context()))
}

func (s Server) renderPage(w http.ResponseWriter, r *http.Request, userID string, lang language.Language) {
	t := texts[lang]
	dp, err := s.Data.userDegreePlan(userID, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, userID, lang)
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
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, userID, lang)
		return
	}
	partialBpBtn := s.BpBtn.PartialComponent(lang)
	partialBpBtnChecked := s.BpBtn.PartialComponentSecond(lang)
	Content(dp, t, partialBpBtn, partialBpBtnChecked).Render(r.Context(), w)
}

func (s Server) pageNotFound(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID, err := s.Auth.UserID(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, "", lang)
		return
	}
	s.Error.RenderPage(w, r, http.StatusNotFound, t.errPageNotFound, t.pageTitle, userID, lang)
}
