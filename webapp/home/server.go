package home

import (
	"net/http"
	"strconv"
	"fmt"

	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/language"
)

//================================================================================
// Server Type
//================================================================================

type Server struct {
	Auth       Authentication
	Error      Error
	Page       Page
	Experiment RecommenderWithAlgoSwitch
	ForYou     Recommender
	Newest     Recommender
	Data       DBManager
	router     http.Handler
}

type Authentication interface {
	// Returns the user ID from an HTTP request.
	UserID(r *http.Request) string
}

type Error interface {
	// Logs the provided error.
	Log(err error)

	// Renders an error message to the user as a floating window, with a status code and localized message.
	Render(w http.ResponseWriter, r *http.Request, code int, userMsg string, lang language.Language)

	// Renders a full error page, including title and user ID, for major errors or page-level failures.
	RenderPage(w http.ResponseWriter, r *http.Request, code int, userMsg string, title string, userID string, lang language.Language)

	// Renders a fallback error page when a regular page cannot be rendered due to an error.
	CannotRenderPage(w http.ResponseWriter, r *http.Request, title string, userID string, err error, lang language.Language)
}

type Page interface {
	// Returns the page view component with injected main content, parameterized by language, title, and user ID.
	// Page adds header with navbar and footer.
	View(main templ.Component, lang language.Language, title string, userID string) templ.Component
}

type Recommender interface {
	Recommend(userID string) ([]string, error)
}

type RecommenderWithAlgoSwitch interface {
	Recommend(userID, student, algoName string, limit int) ([]string, []string, []string, []bool, error)
	Algorithms() ([]string, []bool, error)
	Fit(algo string) error
}

//================================================================================
// Routing
//================================================================================

func (s Server) Router() http.Handler {
	return s.router
}

func (s *Server) Init() {
	router := http.NewServeMux()
	router.HandleFunc("GET /{$}", s.page)
	router.HandleFunc("GET /home/{$}", s.page)
	router.HandleFunc("/", s.pageNotFound)
	router.HandleFunc("GET /recommended/{userID}", s.recommendedPage)
	router.HandleFunc("GET /recommended", s.recommendedPage)
	router.HandleFunc("POST /fit", s.fit)
	s.router = router
}

//================================================================================
// Handlers
//================================================================================

func (s Server) fit(w http.ResponseWriter, r *http.Request) {
	algo := r.FormValue("algo")
	if algo == "" {
		http.Error(w, "algo is required", http.StatusBadRequest)
		return
	}
	err := s.Experiment.Fit(algo)
	if err != nil {
		fmt.Println("error")
		s.Error.Log(errorx.AddContext(err))
	}
	Fit(err != nil).Render(r.Context(), w)
}

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]

	userID := s.Auth.UserID(r)

	var recommended []course
	// recommended, err := s.recommended(userID, lang)
	// if err != nil {
	// 	code, userMsg := errorx.UnwrapError(err, lang)
	// 	s.Error.Log(errorx.AddContext(err))
	// 	s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, userID, lang)
	// 	return
	// }
	newest, err := s.newest(userID, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, userID, lang)
		return
	}
	content := homePage{
		recommendedCourses: recommended,
		newCourses:         newest,
	}

	main := Content(&content, t)
	page := s.Page.View(main, lang, t.pageTitle, userID)
	err = page.Render(r.Context(), w)

	if err != nil {
		s.Error.CannotRenderPage(w, r, t.pageTitle, userID, errorx.AddContext(err), lang)
	}
}

func (s Server) recommendedPage(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]

	userID := s.Auth.UserID(r)
	if testAccount := r.PathValue("userID"); testAccount != "" {
		if testAccount[:5] != "test-" {
			s.Error.RenderPage(w, r, http.StatusBadRequest, "userID must start with 'test-'", t.pageTitle, userID, lang)
			return
		}
		userID = r.PathValue("userID")
	}
	student := r.URL.Query().Get("student")
	algo := r.URL.Query().Get("algo")
	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 50 {
			s.Error.Log(errorx.AddContext(err))
			s.Error.RenderPage(w, r, http.StatusBadRequest, "limit must be between 5 and 50", t.pageTitle, userID, lang)
			return
		}
	}
	var experiment []course
	var finished []course
	var expected []course
	var target []bool
	if len(algo) > 0 {
		var err error
		finished, experiment, expected, target, err = s.experiment(userID, student, algo, limit, lang)
		if err != nil {
			code, userMsg := errorx.UnwrapError(err, lang)
			s.Error.Log(errorx.AddContext(err))
			s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, userID, lang)
			return
		}
	}
	algoSuggestions, algoFit, err := s.Experiment.Algorithms()
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, userID, lang)
		return
	}
	testAccounts, err := s.Data.testAccounts()
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, userID, lang)
		return
	}
	model := recommendedModel{
		student:         student,
		courses:         experiment,
		finished:        finished,
		expected:        expected,
		target:          target,
		algo:            algo,
		algoSuggestions: algoSuggestions,
		algoFit:         algoFit,
		limit:           limit,
		testAccounts:    testAccounts,
	}
	main := Recommended(model, t)
	page := s.Page.View(main, lang, t.pageTitle, userID)
	err = page.Render(r.Context(), w)

	if err != nil {
		s.Error.CannotRenderPage(w, r, t.pageTitle, userID, errorx.AddContext(err), lang)
	}
}

func (s Server) recommended(userID string, lang language.Language) ([]course, error) {
	courses, err := s.ForYou.Recommend(userID)
	if err != nil {
		// TODO: add context
		return nil, err
	}
	similarCourses, err := s.Data.courses(userID, courses, lang)
	if err != nil {
		// TODO: add context
		return nil, err
	}
	return similarCourses, nil
}

func (s Server) newest(userID string, lang language.Language) ([]course, error) {
	courses, err := s.Newest.Recommend(userID)
	if err != nil {
		// TODO: add context
		return nil, err
	}
	newestCourses, err := s.Data.courses(userID, courses, lang)
	if err != nil {
		// TODO: add context
		return nil, err
	}
	return newestCourses, nil
}

func (s Server) experiment(userID, student, algoName string, limit int, lang language.Language) ([]course, []course, []course, []bool, error) {
	finished, recommended, expected, target, err := s.Experiment.Recommend(userID, student, algoName, limit)
	if err != nil {
		// TODO: add context
		return nil, nil, nil, nil, err
	}
	// if len(courses) > 0 {
	finishedCourses, err := s.Data.courses(userID, finished, lang)
	if err != nil {
		// TODO: add context
		return nil, nil, nil, nil, err
	}
	recommendedCourses, err := s.Data.courses(userID, recommended, lang)
	if err != nil {
		// TODO: add context
		return nil, nil, nil, nil, err
	}
	expectedCourses, err := s.Data.courses(userID, expected, lang)
	if err != nil {
		// TODO: add context
		return nil, nil, nil, nil, err
	}
	// }
	return finishedCourses, recommendedCourses, expectedCourses, target, nil
}

func (s Server) pageNotFound(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]

	userID := s.Auth.UserID(r)

	s.Error.RenderPage(w, r, http.StatusNotFound, t.errPageNotFound, t.pageTitle, userID, lang)
}
