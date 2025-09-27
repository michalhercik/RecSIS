package home

import (
	"net/http"

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
	Recommend(userID string, algoName string) ([]string, error)
	Algorithms() ([]string, error)
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
	router.HandleFunc("GET /recommended", s.recommendedPage)
	s.router = router
}

//================================================================================
// Handlers
//================================================================================

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]

	userID := s.Auth.UserID(r)

	recommended, err := s.recommended(userID, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, userID, lang)
		return
	}
	newest, err := s.newest(userID, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, userID, lang)
		return
	}
	//TODO: use constants
	algo := r.URL.Query().Get("algo")
	experiment, err := s.experiment(userID, algo, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, userID, lang)
		return
	}

	content := homePage{
		recommendedCourses: recommended,
		newCourses:         newest,
		experimentCourses:  experiment,
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
	algo := r.URL.Query().Get("algo")
	experiment, err := s.experiment(userID, algo, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, userID, lang)
		return
	}
	algoSuggestions, err := s.Experiment.Algorithms()
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, userID, lang)
		return
	}
	model := recommendedModel{
		courses:         experiment,
		algo:            algo,
		algoSuggestions: algoSuggestions,
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

func (s Server) experiment(userID string, algoName string, lang language.Language) ([]course, error) {
	courses, err := s.Experiment.Recommend(userID, algoName)
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

// func (s Server) fetchCourses(endpoint string, lang language.Language) ([]course, error) {
// 	url := fmt.Sprintf("%s/%s?lang=%s", s.Recommender, endpoint, lang)
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, errorx.NewHTTPErr(
// 			errorx.AddContext(err, errorx.P("URL", url)),
// 			http.StatusServiceUnavailable,
// 			texts[lang].errRecommenderUnavailable,
// 		)
// 	}

// 	defer resp.Body.Close()
// 	if resp.StatusCode != http.StatusOK {
// 		return nil, errorx.NewHTTPErr(
// 			errorx.AddContext(fmt.Errorf("unexpected status code: %d", resp.StatusCode), errorx.P("URL", url)),
// 			resp.StatusCode,
// 			texts[lang].errRecommenderUnavailable,
// 		)
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, errorx.NewHTTPErr(
// 			errorx.AddContext(fmt.Errorf("failed to read response body: %w", err), errorx.P("URL", url)),
// 			http.StatusServiceUnavailable,
// 			texts[lang].errCannotLoadCourses,
// 		)
// 	}

// 	var courses []course
// 	err = json.Unmarshal(body, &courses)
// 	if err != nil {
// 		return nil, errorx.NewHTTPErr(
// 			errorx.AddContext(fmt.Errorf("failed to unmarshal response: %w", err), errorx.P("URL", url)),
// 			http.StatusInternalServerError,
// 			texts[lang].errCannotLoadCourses,
// 		)
// 	}

// 	return courses, nil
// }

func (s Server) pageNotFound(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]

	userID := s.Auth.UserID(r)

	s.Error.RenderPage(w, r, http.StatusNotFound, t.errPageNotFound, t.pageTitle, userID, lang)
}
