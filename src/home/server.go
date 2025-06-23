package home

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/language"
)

//================================================================================
// Server Type
//================================================================================

type Server struct {
	Auth        Authentication
	Error       Error
	Page        Page
	Recommender string
	router      http.Handler
}

func (s Server) Router() http.Handler {
	return s.router
}

type Authentication interface {
	UserID(r *http.Request) (string, error)
}

type Error interface {
	Log(err error)
	Render(w http.ResponseWriter, r *http.Request, code int, userMsg string, lang language.Language)
	RenderPage(w http.ResponseWriter, r *http.Request, code int, userMsg string, title string, userID string, lang language.Language)
	CannotRenderPage(w http.ResponseWriter, r *http.Request, title string, userID string, err error, lang language.Language)
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
	router.HandleFunc("GET /home/", s.page)

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

	recommended, err := s.fetchCourses("recommended", lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, userID, lang)
		return
	}

	newest, err := s.fetchCourses("newest", lang)
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
	err = s.Page.View(main, lang, t.pageTitle, userID).Render(r.Context(), w)

	if err != nil {
		s.Error.CannotRenderPage(w, r, t.pageTitle, userID, errorx.AddContext(err), lang)
	}
}

func (s Server) fetchCourses(endpoint string, lang language.Language) ([]course, error) {
	url := fmt.Sprintf("%s/%s?lang=%s", s.Recommender, endpoint, lang)
	resp, err := http.Get(url)
	if err != nil {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(err, errorx.P("URL", url)),
			http.StatusServiceUnavailable,
			texts[lang].errRecommenderUnavailable,
		)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("unexpected status code: %d", resp.StatusCode), errorx.P("URL", url)),
			resp.StatusCode,
			texts[lang].errRecommenderUnavailable,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("failed to read response body: %w", err), errorx.P("URL", url)),
			http.StatusServiceUnavailable,
			texts[lang].errCannotLoadCourses,
		)
	}

	var courses []course
	err = json.Unmarshal(body, &courses)
	if err != nil {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("failed to unmarshal response: %w", err), errorx.P("URL", url)),
			http.StatusInternalServerError,
			texts[lang].errCannotLoadCourses,
		)
	}

	return courses, err
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
