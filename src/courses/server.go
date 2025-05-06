package courses

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/michalhercik/RecSIS/courses/internal/filter"
	"github.com/michalhercik/RecSIS/dbcourse"
	"github.com/michalhercik/RecSIS/language"
)

const courseIndex = "courses"
const pageParam = "page"
const hitsPerPageParam = "hitsPerPage"

type Server struct {
	router *http.ServeMux
	Data   DataManager
	// SearchParam string
	Search  SearchEngine
	Auth    Authentication
	filters filter.Filters
	BpBtn   BlueprintAddButton
	Page    Page
	// PageTempl   func(templ.Component, language.Language, string) templ.Component
}

//================================================================================
// Interface
//================================================================================

func (s *Server) Init() {
	s.initFilters()
	s.initRouter()
}

func (s Server) Router() http.Handler {
	return s.router
}

//================================================================================
// Init
//================================================================================

func (s *Server) initRouter() {
	router := http.NewServeMux()
	router.HandleFunc("GET /", s.page)
	router.HandleFunc("GET /search", s.content)
	// router.HandleFunc("GET /quicksearch", s.quickSearch)
	router.HandleFunc("POST /blueprint/{coursecode}", s.addCourseToBlueprint)
	s.router = router
}

func (s *Server) initFilters() {
	filters, err := s.Data.Filters()
	if err != nil {
		log.Fatalf("failed to load filters: %v", err)
	}
	s.filters = filters
}

//================================================================================
// Handlers
//================================================================================

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
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
	numberOfBlueprintYears, err := s.BpBtn.NumberOfYears(req.userID)
	if err != nil {
		log.Printf("numberOfBlueprintYears: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result.templ = s.BpBtn.PartialComponent(numberOfBlueprintYears, lang)
	main := Content(&result, t)
	s.Page.View(main, lang, t.Title, req.query).Render(r.Context(), w)
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

	numberOfBlueprintYears, err := s.BpBtn.NumberOfYears(req.userID)
	if err != nil {
		log.Printf("numberOfBlueprintYears: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	coursesPage.templ = s.BpBtn.PartialComponent(numberOfBlueprintYears, lang)
	Content(&coursesPage, t).Render(r.Context(), w)
}

func (s Server) quickSearch(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	query := r.FormValue(s.Page.SearchParam())
	req := QuickRequest{
		query:    query,
		indexUID: courseIndex,
		limit:    5,
		offset:   0,
		lang:     lang,
	}
	res, err := s.Search.QuickSearch(req)
	if err != nil {
		log.Printf("quickSearch: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	QuickResults(&res, t).Render(r.Context(), w)
}

func (s Server) parseQueryRequest(w http.ResponseWriter, r *http.Request) (Request, error) {
	var req Request
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
	filter, err := s.filters.ParseURLQuery(r.URL.Query())
	if err != nil {
		// TODO: handle error
		log.Printf("search error: %v", err)
	}

	req = Request{
		userID:      userID,
		query:       query,
		indexUID:    courseIndex,
		page:        page,
		hitsPerPage: hitsPerPage,
		lang:        lang,
		filter:      filter,
		facets:      s.filters.Facets,
	}
	return req, nil
}

func (s Server) search(req Request, httpReq *http.Request) (coursesPage, error) {
	// search for courses
	var result coursesPage
	searchResponse, err := s.Search.Search(req)
	if err != nil {
		return result, err
	}
	coursesData, err := s.Data.Courses(req.userID, searchResponse.Courses, req.lang)
	if err != nil {
		return result, err
	}
	// paramLabels, err := s.Data.ParamLabels(req.lang)
	// if err != nil {
	// 	return result, err
	// }
	// facets := filter.MakeFacetDistribution(searchResponse.FacetDistribution, paramLabels)
	result = coursesPage{
		courses:     coursesData,
		page:        int(req.page),
		pageSize:    int(req.hitsPerPage),
		searchParam: s.Page.SearchParam(),
		totalPages:  searchResponse.TotalPages,
		search:      req.query,
		facets:      filter.IterFiltersWithFacets(s.filters, searchResponse.FacetDistribution, httpReq.URL.Query(), req.lang),
	}
	return result, nil
}

func (s Server) facetDistribution(lang language.Language) (coursesPage, error) {
	var result coursesPage
	// f, err := s.Search.FacetDistribution()
	// if err != nil {
	// 	return result, err
	// }
	// param, err := s.Data.ParamLabels(lang)
	// if err != nil {
	// 	return result, err
	// }
	// facets := filter.MakeFacetDistribution(f, param)
	// result = coursesPage{
	// 	facets: facets,
	// }
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

	return t.Utils.LangLink("/courses?" + queryValues.Encode())
}

func (s Server) addCourseToBlueprint(w http.ResponseWriter, r *http.Request) {
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Printf("addCourseToBlueprint: %v", err)
		return
	}
	lang := language.FromContext(r.Context())
	courseCode := r.PathValue("coursecode")
	year, semester, err := parseYearSemester(r)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		log.Printf("addCourseToBlueprint: %v", err)
		return
	}
	_, err = s.BpBtn.Action(userID, courseCode, year, semester)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("failed to create button: %v", err)
		return
	}

	t := texts[lang]
	courses, err := s.Data.Courses(userID, []string{courseCode}, lang)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("failed to create button: %v", err)
		return
	}
	numberOfYears, err := s.BpBtn.NumberOfYears(userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("failed to create button: %v", err)
		return
	}
	btn := s.BpBtn.PartialComponent(numberOfYears, lang)
	CourseCard(&courses[0], t, btn).Render(r.Context(), w)

	// btn.Render(r.Context(), w)
}

func parseYearSemester(r *http.Request) (int, dbcourse.SemesterAssignment, error) {
	year, err := parseYear(r)
	if err != nil {
		return 0, 0, err
	}
	semester, err := parseSemester(r)
	if err != nil {
		return year, 0, err
	}
	return year, semester, nil
}

func parseYear(r *http.Request) (int, error) {
	year, err := strconv.Atoi(r.FormValue("year"))
	if err != nil {
		return year, err
	}
	return year, nil
}

func parseSemester(r *http.Request) (dbcourse.SemesterAssignment, error) {
	semesterInt, err := strconv.Atoi(r.FormValue("semester"))
	if err != nil {
		return 0, err
	}
	semester := dbcourse.SemesterAssignment(semesterInt)
	return semester, nil
}
