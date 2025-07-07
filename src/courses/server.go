package courses

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/filters"
	"github.com/michalhercik/RecSIS/language"
)

//================================================================================
// Server Type
//================================================================================

type Server struct {
	Auth    Authentication
	BpBtn   BlueprintAddButton
	Data    DBManager
	Error   Error
	Filters filters.Filters
	Page    Page
	router  http.Handler
	Search  searchEngine
}

func (s *Server) Init() {
	if err := s.Filters.Init(); err != nil {
		log.Fatal("courses.Init: ", err)
	}
	s.initRouter()
}

func (s Server) Router() http.Handler {
	return s.router
}

type Authentication interface {
	UserID(r *http.Request) string
}

type BlueprintAddButton interface {
	PartialComponent(lang language.Language) PartialBlueprintAdd
	ParseRequest(r *http.Request, additionalCourses []string) ([]string, int, int, error)
	Action(userID string, year int, semester int, lang language.Language, course ...string) ([]int, error)
}

type PartialBlueprintAdd = func(hxSwap, hxTarget, hxInclude string, years []bool, course string) templ.Component

type Error interface {
	Log(err error)
	Render(w http.ResponseWriter, r *http.Request, code int, userMsg string, lang language.Language)
	RenderPage(w http.ResponseWriter, r *http.Request, code int, userMsg string, title string, userID string, lang language.Language)
	CannotRenderPage(w http.ResponseWriter, r *http.Request, title string, userID string, err error, lang language.Language)
	CannotRenderComponent(w http.ResponseWriter, r *http.Request, err error, lang language.Language)
}

type Page interface {
	View(main templ.Component, lang language.Language, title string, searchParam string, userID string) templ.Component
	SearchParam() string
}

//================================================================================
// Routing
//================================================================================

func (s *Server) initRouter() {
	router := http.NewServeMux()
	router.HandleFunc("GET /{$}", s.page)
	router.HandleFunc("GET /search", s.content)
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
	userID := s.Auth.UserID(r)
	req, err := s.parseQueryRequest(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		userMsg = fmt.Sprintf("%s: %s", t.errCannotSearchCourses, userMsg)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, userID, lang)
		return
	}
	result, err := s.search(req, r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, userID, lang)
		return
	}
	result.bpBtn = s.BpBtn.PartialComponent(lang)
	main := Content(&result, t)
	err = s.Page.View(main, lang, t.pageTitle, req.query, userID).Render(r.Context(), w)
	if err != nil {
		s.Error.CannotRenderPage(w, r, t.pageTitle, userID, errorx.AddContext(err), lang)
	}
}

func (s Server) content(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	req, err := s.parseQueryRequest(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		userMsg = fmt.Sprintf("%s: %s", t.errCannotSearchCourses, userMsg)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	coursesPage, err := s.search(req, r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}

	// Set the HX-Push-Url header to update the browser URL without a full reload
	w.Header().Set("HX-Push-Url", s.parseUrl(r.URL.Query(), lang))

	coursesPage.bpBtn = s.BpBtn.PartialComponent(lang)
	err = Content(&coursesPage, t).Render(r.Context(), w)
	if err != nil {
		s.Error.CannotRenderPage(w, r, t.pageTitle, "", errorx.AddContext(err), lang)
	}
}

func (s Server) parseQueryRequest(r *http.Request) (request, error) {
	var req request
	var err error
	userID := s.Auth.UserID(r)
	lang := language.FromContext(r.Context())
	query := r.FormValue(s.Page.SearchParam())
	pageString := r.FormValue(pageParam)
	page := 1 // default to page 1 if not specified
	if pageString != "" {
		page, err = strconv.Atoi(pageString)
		if err != nil || page < 1 {
			return req, errorx.NewHTTPErr(
				errorx.AddContext(fmt.Errorf("invalid page number: '%s'", pageString)),
				http.StatusBadRequest,
				texts[lang].errInvalidPageNumber,
			)
		}
	}
	hitsPerPageString := r.FormValue(hitsPerPageParam)
	hitsPerPage := coursesPerPage // default to coursesPerPage if not specified
	if hitsPerPageString != "" {
		hitsPerPage, err = strconv.Atoi(hitsPerPageString)
		if err != nil || hitsPerPage < 1 {
			return req, errorx.NewHTTPErr(
				errorx.AddContext(fmt.Errorf("invalid number of courses per page: %d", hitsPerPage)),
				http.StatusBadRequest,
				texts[lang].errInvalidNumberOfCourses,
			)
		}
	}
	filter, err := s.Filters.ParseURLQuery(r.URL.Query(), lang)
	if err != nil {
		return req, errorx.AddContext(err)
	}

	req = request{
		userID:      userID,
		query:       query,
		indexUID:    courseIndex,
		page:        page,
		hitsPerPage: hitsPerPage,
		lang:        lang,
		filter:      filter,
		facets:      s.Filters.Facets(),
	}
	return req, nil
}

func (s Server) search(req request, httpReq *http.Request) (coursesPage, error) {
	// search for courses
	var result coursesPage
	searchResponse, err := s.Search.Search(req)
	if err != nil {
		return result, errorx.AddContext(err)
	}
	coursesData, err := s.Data.courses(req.userID, searchResponse.Courses, req.lang)
	if err != nil {
		return result, errorx.AddContext(err)
	}
	result = coursesPage{
		courses:     coursesData,
		page:        int(req.page),
		pageSize:    int(req.hitsPerPage),
		searchParam: s.Page.SearchParam(),
		totalPages:  searchResponse.TotalPages,
		search:      req.query,
		facets:      s.Filters.IterFiltersWithFacets(searchResponse.FacetDistribution, httpReq.URL.Query(), req.lang),
	}
	return result, nil
}

func (s Server) parseUrl(queryValues url.Values, lang language.Language) string {
	// exclude default values from the URL
	if queryValues.Get(s.Page.SearchParam()) == "" {
		queryValues.Del(s.Page.SearchParam())
	}
	if queryValues.Get(pageParam) == "1" {
		queryValues.Del(pageParam)
	}
	// TODO: possibly add more defaults to exclude

	return fmt.Sprintf("%s?%s", lang.LocalizeURL("/courses/"), queryValues.Encode())
}

func (s Server) addCourseToBlueprint(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID := s.Auth.UserID(r)
	courseCodes, year, semester, err := s.BpBtn.ParseRequest(r, nil)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	if len(courseCodes) != 1 {
		s.Error.Log(errorx.AddContext(fmt.Errorf("expected exactly one course code, got %d", len(courseCodes)), errorx.P("courseCodes", courseCodes)))
		s.Error.Render(w, r, http.StatusBadRequest, t.errUnexpectedNumberOfCourses, lang)
		return
	}
	courseCode := courseCodes[0]
	_, err = s.BpBtn.Action(userID, year, semester, lang, courseCode)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}

	course, err := s.Data.courses(userID, []string{courseCode}, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	btn := s.BpBtn.PartialComponent(lang)
	err = CourseCard(&course[0], t, btn).Render(r.Context(), w)
	if err != nil {
		s.Error.CannotRenderComponent(w, r, errorx.AddContext(err), lang)
	}
}

func (s Server) pageNotFound(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID := s.Auth.UserID(r)
	s.Error.RenderPage(w, r, http.StatusNotFound, t.errPageNotFound, t.pageTitle, userID, lang)
}
