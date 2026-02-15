package degreeplans

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/degreeplans/compare"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/filters"
	"github.com/michalhercik/RecSIS/language"
)

//================================================================================
// Server Type
//================================================================================

type Server struct {
	Auth          Authentication
	Data          DBManager
	Error         Error
	Filters       filters.Filters
	Page          Page
	router        http.Handler
	Search        searchEngine
	compareServer *compare.Server
}

func (s *Server) Init() {
	if err := s.Filters.Init(); err != nil {
		log.Fatal("degreeplan.Init: ", err)
	}
	s.initCompareServer()
	s.initRouter()
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

	// Renders a floating window with error when any component cannot be rendered due to an error.
	CannotRenderComponent(w http.ResponseWriter, r *http.Request, err error, lang language.Language)
}

type Page interface {
	// Returns the page view component with injected main content, parameterized by language, title, and user ID.
	// Page adds header with navbar and footer.
	View(main templ.Component, lang language.Language, title string, userID string) templ.Component
}

//================================================================================
// Subserver Initialization
//================================================================================

func (s *Server) initCompareServer() {
	s.compareServer = &compare.Server{
		Auth:  s.Auth,
		Data:  compare.DBManager{DB: s.Data.DB},
		Error: s.Error,
		Page:  s.Page,
	}
	s.compareServer.Init()
}

//================================================================================
// Routing
//================================================================================

func (s Server) Router() http.Handler {
	return s.router
}

func (s *Server) initRouter() {
	router := http.NewServeMux()
	router.HandleFunc("GET /{$}", s.searchPage)
	router.HandleFunc("GET /search", s.searchContent)
	s.handleCompareRoutes(router)
	router.HandleFunc("/", s.pageNotFound)
	s.router = router
}

func (s *Server) handleCompareRoutes(router *http.ServeMux) {
	router.Handle(comparePrefix, http.StripPrefix(comparePrefix[:len(comparePrefix)-1], s.compareServer.Router()))
}

//================================================================================
// Handlers
//================================================================================

func (s Server) searchPage(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID := s.Auth.UserID(r)
	degreePlanSearchContent, err := s.getDegreePlanData(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, userID, lang)
		return
	}
	main := Content(degreePlanSearchContent, t)
	page := s.Page.View(main, lang, t.pageTitle, userID)
	err = page.Render(r.Context(), w)
	if err != nil {
		s.Error.CannotRenderPage(w, r, t.pageTitle, userID, errorx.AddContext(err), lang)
	}
}

func (s Server) searchContent(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	degreePlanSearchContent, err := s.getDegreePlanData(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	w.Header().Set("HX-Push-Url", s.parseUrl(r.URL.Query(), lang))
	content := FilterResults(degreePlanSearchContent, t)
	err = content.Render(r.Context(), w)
	if err != nil {
		s.Error.CannotRenderComponent(w, r, errorx.AddContext(err), lang)
	}
}

func (s Server) getDegreePlanData(r *http.Request) (*degreePlanSearchPage, error) {
	req, err := s.parseRequest(r)
	if err != nil {
		return nil, errorx.AddContext(err)
	}
	res, err := s.Search.Search(req)
	if err != nil {
		return nil, errorx.AddContext(err)
	}
	result, err := s.buildDegreePlanPage(req, res, r)
	if err != nil {
		return nil, errorx.AddContext(err)
	}
	return &result, nil
}

func (s Server) parseRequest(r *http.Request) (request, error) {
	var req request
	var err error
	userID := s.Auth.UserID(r)
	lang := language.FromContext(r.Context())
	query := r.FormValue(searchDegreePlanName)
	filter, err := s.Filters.ParseURLQuery(r.URL.Query(), lang)
	if err != nil {
		return req, errorx.AddContext(err)
	}
	req = request{
		userID:   userID,
		query:    query,
		indexUID: SearchIndex,
		lang:     lang,
		filter:   filter,
		facets:   s.Filters.Facets(),
	}
	return req, nil
}

func (s Server) buildDegreePlanPage(req request, res response, httpReq *http.Request) (degreePlanSearchPage, error) {
	var result degreePlanSearchPage
	degreePlanMetadata, err := s.Data.degreePlanMetadata(res.DegreePlanCodes, req.lang)
	if err != nil {
		return result, errorx.AddContext(err)
	}
	compareCode := httpReq.FormValue(CompareUrlParam)
	result = degreePlanSearchPage{
		filters:     s.Filters.FiltersMapWithFacets(res.FacetDistribution, httpReq.URL.Query(), req.lang),
		results:     degreePlanMetadata,
		searchQuery: req.query,
		selectedPlan: selectedPlan{
			isAnySelected: compareCode != "",
			code:          compareCode,
		},
	}
	return result, nil
}

func (s Server) pageNotFound(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID := s.Auth.UserID(r)
	s.Error.RenderPage(w, r, http.StatusNotFound, t.errPageNotFound, t.pageTitle, userID, lang)
}

func (s Server) parseUrl(queryValues url.Values, lang language.Language) string {
	// exclude default values from the URL
	if queryValues.Get(searchDegreePlanName) == "" {
		queryValues.Del(searchDegreePlanName)
	}
	if queryValues.Get(CompareUrlParam) == "" {
		queryValues.Del(CompareUrlParam)
	}
	// build URL
	url := lang.LocalizeURL("/degreeplans/")
	if queryValues.Encode() != "" {
		url = fmt.Sprintf("%s?%s", url, queryValues.Encode())
	}
	return url
}
