package courses

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/a-h/templ"
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
	Filters filters.Filters
	Page    Page
	router  *http.ServeMux
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
	UserID(r *http.Request) (string, error)
}

type BlueprintAddButton interface {
	PartialComponent(lang language.Language) PartialBlueprintAdd
	ParseRequest(r *http.Request) ([]string, int, int, error)
	Action(userID string, year int, semester int, course ...string) ([]int, error)
}

type PartialBlueprintAdd = func(hxSwap, hxTarget, hxInclude string, years []bool, course ...string) templ.Component

type Page interface {
	View(main templ.Component, lang language.Language, title string, searchParam string, userID string) templ.Component
	SearchParam() string
}

//================================================================================
// Routing
//================================================================================

func (s *Server) initRouter() {
	router := http.NewServeMux()
	router.HandleFunc("GET /", s.page)
	router.HandleFunc("GET /search", s.content)
	router.HandleFunc("POST /blueprint", s.addCourseToBlueprint)
	s.router = router
}

//================================================================================
// Handlers
//================================================================================

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Printf("page: %v", err)
		return
	}
	req, err := s.parseQueryRequest(w, r)
	if err != nil {
		// TODO: handle error
		log.Printf("search: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, err := s.search(req, r)
	if err != nil {
		log.Printf("search: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result.bpBtn = s.BpBtn.PartialComponent(lang)
	main := Content(&result, t)
	s.Page.View(main, lang, t.pageTitle, req.query, userID).Render(r.Context(), w)
}

func (s Server) content(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	req, err := s.parseQueryRequest(w, r)
	if err != nil {
		// TODO: handle error
		log.Printf("search: %v", err)
		return
	}
	coursesPage, err := s.search(req, r)
	if err != nil {
		log.Printf("search: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the HX-Push-Url header to update the browser URL without a full reload
	w.Header().Set("HX-Push-Url", s.parseUrl(r.URL.Query(), t))

	coursesPage.bpBtn = s.BpBtn.PartialComponent(lang)
	Content(&coursesPage, t).Render(r.Context(), w)
}

func (s Server) parseQueryRequest(w http.ResponseWriter, r *http.Request) (request, error) {
	var req request
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return req, err
	}
	lang := language.FromContext(r.Context())
	query := r.FormValue(s.Page.SearchParam())
	page, err := strconv.Atoi(r.FormValue(pageParam))
	if err != nil {
		page = 1
	}
	hitsPerPage, err := strconv.Atoi(r.FormValue(hitsPerPageParam))
	if err != nil {
		hitsPerPage = coursesPerPage
	}
	filter, err := s.Filters.ParseURLQuery(r.URL.Query())
	if err != nil {
		// TODO: handle error
		log.Printf("search error: %v", err)
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
		return result, err
	}
	coursesData, err := s.Data.courses(req.userID, searchResponse.Courses, req.lang)
	if err != nil {
		return result, err
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

func (s Server) parseUrl(queryValues url.Values, t text) string {
	// exclude default values from the URL
	if queryValues.Get(s.Page.SearchParam()) == "" {
		queryValues.Del(s.Page.SearchParam())
	}
	if queryValues.Get(pageParam) == "1" {
		queryValues.Del(pageParam)
	}
	// TODO: possibly add more defaults to exclude

	return fmt.Sprintf("%s?%s", t.language.LocalizeURL("/courses"), queryValues.Encode())
}

func (s Server) addCourseToBlueprint(w http.ResponseWriter, r *http.Request) {
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Printf("addCourseToBlueprint: %v", err)
		return
	}
	lang := language.FromContext(r.Context())
	courseCodes, year, semester, err := s.BpBtn.ParseRequest(r)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		log.Printf("addCourseToBlueprint: %v", err)
		return
	}
	_, err = s.BpBtn.Action(userID, year, semester, courseCodes[0])
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("failed to create button: %v", err)
		return
	}

	t := texts[lang]
	courses, err := s.Data.courses(userID, courseCodes, lang)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("failed to create button: %v", err)
		return
	}
	btn := s.BpBtn.PartialComponent(lang)
	CourseCard(&courses[0], t, btn).Render(r.Context(), w)

	// btn.Render(r.Context(), w)
}
