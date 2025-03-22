package courses

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/michalhercik/RecSIS/courses/internal/filter"
)

const courseIndex = "courses"

type Server struct {
	Data   DataManager
	Search SearchEngine
}

func (s Server) Register(router *http.ServeMux, prefix string) {
	//router.HandleFunc(fmt.Sprintf("GET %s", prefix), s.page) // TODO get language from http header
	router.HandleFunc(fmt.Sprintf("GET /cs%s", prefix), s.csPage)
	router.HandleFunc(fmt.Sprintf("GET /en%s", prefix), s.enPage)
	//router.HandleFunc(fmt.Sprintf("GET %s/search", prefix), s.content) // TODO get language from http header
	router.HandleFunc(fmt.Sprintf("GET /cs%s/search", prefix), s.csContent)
	router.HandleFunc(fmt.Sprintf("GET /en%s/search", prefix), s.enContent)
	//router.HandleFunc(fmt.Sprintf("GET %s/quicksearch", prefix), s.quickSearch) // TODO get language from http header
	router.HandleFunc(fmt.Sprintf("GET /cs%s/quicksearch", prefix), s.csQuickSearch)
	router.HandleFunc(fmt.Sprintf("GET /en%s/quicksearch", prefix), s.enQuickSearch)
}

func (s Server) csPage(w http.ResponseWriter, r *http.Request) {
	s.page(w, r, cs, texts["cs"])
}

func (s Server) enPage(w http.ResponseWriter, r *http.Request) {
	s.page(w, r, en, texts["en"])
}

func (s Server) page(w http.ResponseWriter, r *http.Request, lang Language, t text) {
	req, err := parseQueryRequest(r, lang)
	if err != nil {
		// TODO: handle error
		log.Printf("search: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := s.search(req)
	if err != nil {
		log.Printf("search: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// coursesPage := createPageContent(res, req)
	coursesPage := res
	Page(&coursesPage, t).Render(r.Context(), w)
}

func (s Server) csContent(w http.ResponseWriter, r *http.Request) {
	s.content(w, r, cs, texts["cs"])
}

func (s Server) enContent(w http.ResponseWriter, r *http.Request) {
	s.content(w, r, en, texts["en"])
}

func (s Server) content(w http.ResponseWriter, r *http.Request, lang Language, t text) {
	req, err := parseQueryRequest(r, lang)
	if err != nil {
		// TODO: handle error
		log.Printf("search: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := s.search(req)
	if err != nil {
		log.Printf("search: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// coursesPage := createPageContent(res, req)
	coursesPage := res
	Courses(&coursesPage, t).Render(r.Context(), w)
}

func (s Server) csQuickSearch(w http.ResponseWriter, r *http.Request) {
	s.quickSearch(w, r, cs, texts["cs"])
}

func (s Server) enQuickSearch(w http.ResponseWriter, r *http.Request) {
	s.quickSearch(w, r, en, texts["en"])
}

func (s Server) quickSearch(w http.ResponseWriter, r *http.Request, lang Language, t text) {
	query := r.FormValue("search")
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

func parseQueryRequest(r *http.Request, lang Language) (Request, error) {
	var req Request
	sessionCookie, err := r.Cookie("recsis_session_key")
	if err != nil {
		return req, err
	}
	query := r.FormValue("search")
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		page = 1
	}
	hitsPerPage, err := strconv.Atoi(r.FormValue("hitsPerPage"))
	if err != nil {
		hitsPerPage = coursesPerPage
	}

	// TODO change language based on URL
	req = Request{
		sessionID:   sessionCookie.Value,
		query:       query,
		indexUID:    courseIndex,
		page:        page,
		hitsPerPage: hitsPerPage,
		lang:        lang,
	}
	return req, nil
}

func (s Server) search(req Request) (coursesPage, error) {
	// search for courses
	var result coursesPage
	searchResponse, err := s.Search.Search(req)
	if err != nil {
		return result, err
	}
	coursesData, err := s.Data.Courses(req.sessionID, searchResponse.Courses, req.lang)
	if err != nil {
		return result, err
	}
	paramLabels, err := s.Data.ParamLabels(req.lang)
	if err != nil {
		return result, err
	}
	facets := filter.MakeFacetDistribution(searchResponse.FacetDistribution, paramLabels)
	result = coursesPage{
		courses:    coursesData,
		page:       int(req.page),
		pageSize:   int(req.hitsPerPage),
		totalPages: searchResponse.TotalPages,
		search:     req.query,
		facets:     facets,
	}
	return result, nil
}
