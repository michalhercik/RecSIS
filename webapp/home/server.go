package home

import (
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/language"
	"github.com/michalhercik/RecSIS/recommend"
)

//================================================================================
// Server Type
//================================================================================

type Server struct {
	Auth   Authentication
	Error  Error
	Page   Page
	ForYou recommend.ForYouRecommender
	Newest Recommender
	Data   DBManager
	router http.Handler
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
	router.HandleFunc("GET /foryou", s.viewAll)
	router.HandleFunc("GET /foryou/content", s.forYouContent)
	router.HandleFunc("/", s.pageNotFound)
	s.router = router
}

//================================================================================
// Handlers
//================================================================================

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]

	userID := s.Auth.UserID(r)

	forYou, categories, err := s.forYou(userID, 0, 20, lang)
	if err != nil {
		code, _ := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, t.errForYou, t.pageTitle, userID, lang)
		return
	}
	content := homePage{
		forYou: forYou,
		categories: categories,
	}

	main := Content(&content, t)
	page := s.Page.View(main, lang, t.pageTitle, userID)
	err = page.Render(r.Context(), w)

	if err != nil {
		s.Error.CannotRenderPage(w, r, t.pageTitle, userID, errorx.AddContext(err), lang)
	}
}

func (s Server) viewAll(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID := s.Auth.UserID(r)

	forYou, _, err := s.forYou(userID, 0, 20, lang)
	if err != nil {
		code, _ := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, t.errForYou, t.pageTitle, userID, lang)
		return
	}

	content := forYouPage{
		forYou: forYou,
	}
	main := ForYouPage(content, t)
	page := s.Page.View(main, lang, t.pageTitle, userID)
	err = page.Render(r.Context(), w)

	if err != nil {
		s.Error.CannotRenderPage(w, r, t.pageTitle, userID, errorx.AddContext(err), lang)
	}
}

func (s Server) forYou(userID string, offset, limit int, lang language.Language) ([]course, []category, error) {
	res, err := s.ForYou.Recommend(userID, offset, limit)
	if err != nil {
		return nil, nil, errorx.AddContext(err)
	}
	allCourses, err := s.Data.courses(userID, res.All, lang)
	if err != nil {
		return nil, nil, errorx.AddContext(err)
	}
	coursesMap := make(map[string]course, len(allCourses))
	for _, course := range allCourses {
		coursesMap[course.Code] = course
	}
	forYou := make([]course, len(res.Courses))
	for i, courseCode := range res.Courses {
		forYou[i] = coursesMap[courseCode]
	}
	categories := make([]category, len(res.Categories.Names))
	for i, name := range res.Categories.Names {
		catCourses := make([]course, len(res.Categories.Values[i]))
		for i, courseCode := range res.Categories.Values[i] {
			catCourses[i] = coursesMap[courseCode]
		}
		categories[i] = category{
			name:   name,
			courses: catCourses,
		}
	}
	return forYou, categories, nil
}

func (s Server) forYouContent(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID := s.Auth.UserID(r)

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}
	limit := 20

	forYou, _, err := s.forYou(userID, offset, limit, lang)
	if err != nil {
		code, _ := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, t.errForYou, t.pageTitle, userID, lang)
		return
	}

	content := forYouPage{
		forYou: forYou,
	}
	err = ForYouContent(content, t).Render(r.Context(), w)
	if err != nil {
		s.Error.CannotRenderPage(w, r, t.pageTitle, userID, errorx.AddContext(err), lang)
	}
}

func (s Server) pageNotFound(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]

	userID := s.Auth.UserID(r)

	s.Error.RenderPage(w, r, http.StatusNotFound, t.errPageNotFound, t.pageTitle, userID, lang)
}
