package courses

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/michalhercik/RecSIS/courses/internal/filter"
	"github.com/michalhercik/RecSIS/language"
)

const courseIndex = "courses"
const searchParam = "search"
const pageParam = "page"
const hitsPerPageParam = "hitsPerPage"

type Server struct {
	Data   DataManager
	Search SearchEngine
}

func (s Server) Register(router *http.ServeMux, prefix string) {
	lr := language.LanguageRouter{Router: router}
	lr.HandleLangFunc(prefix, http.MethodGet, s.page)
	lr.HandleLangFunc(prefix+"/search/", http.MethodGet, s.content)
	lr.HandleLangFunc(prefix+"/quicksearch/", http.MethodGet, s.quickSearch)
}

func (s Server) page(w http.ResponseWriter, r *http.Request, lang language.Language) {
	t := texts[lang]
	req, err := parseQueryRequest(r, lang)
	if err != nil {
		// TODO: handle error
		log.Printf("search: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, err := s.search(req)
	if err != nil {
		log.Printf("search: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	Page(&result, t).Render(r.Context(), w)
}

func (s Server) content(w http.ResponseWriter, r *http.Request, lang language.Language) {
	t := texts[lang]
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

	// Set the HX-Push-Url header to update the browser URL without a full reload
	w.Header().Set("HX-Push-Url", parseUrl(r.URL.Query(), t))

	// coursesPage := createPageContent(res, req)
	coursesPage := res
	// TODO: return page
	Content(&coursesPage, t).Render(r.Context(), w)
}

func (s Server) quickSearch(w http.ResponseWriter, r *http.Request, lang language.Language) {
	t := texts[lang]
	query := r.FormValue(searchParam)
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

func parseQueryRequest(r *http.Request, lang language.Language) (Request, error) {
	var req Request
	sessionCookie, err := r.Cookie("recsis_session_key")
	if err != nil {
		return req, err
	}
	query := r.FormValue(searchParam)
	page, err := strconv.Atoi(r.FormValue(pageParam))
	if err != nil {
		page = 1
	}
	hitsPerPage, err := strconv.Atoi(r.FormValue(hitsPerPageParam))
	if err != nil {
		hitsPerPage = coursesPerPage
	}
	filter, err := filter.ParseFilters(r.URL.Query())
	if err != nil {
		// TODO: handle error
		log.Printf("search error: %v", err)
	}

	req = Request{
		sessionID:   sessionCookie.Value,
		query:       query,
		indexUID:    courseIndex,
		page:        page,
		hitsPerPage: hitsPerPage,
		lang:        lang,
		filter:      filter,
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

func (s Server) facetDistribution(lang language.Language) (coursesPage, error) {
	var result coursesPage
	f, err := s.Search.FacetDistribution()
	if err != nil {
		return result, err
	}
	param, err := s.Data.ParamLabels(lang)
	if err != nil {
		return result, err
	}
	facets := filter.MakeFacetDistribution(f, param)
	result = coursesPage{
		facets: facets,
	}
	return result, nil
}

func parseUrl(queryValues url.Values, t text) string {
	// exclude default values from the URL
	if queryValues.Get("search") == "" {
		queryValues.Del("search")
	}
	if queryValues.Get("page") == "1" {
		queryValues.Del("page")
	}
	// TODO: possibly add more defaults to exclude

	return t.Utils.LangLink("/courses?" + queryValues.Encode())
}
