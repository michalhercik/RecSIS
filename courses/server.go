package courses

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	coursesPage := createPageContent(res, req)
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
	coursesPage := createPageContent(res, req)
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
	res, err := s.Search.QuickSearch(&req)
	if err != nil {
		log.Printf("quickSearch: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	QuickResults(res, t).Render(r.Context(), w)
}

func parseQueryRequest(r *http.Request, lang Language) (*Request, error) {
	sessionCookie, err := r.Cookie("recsis_session_key")
	if err != nil {
		return nil, err
	}
	query := r.FormValue("search")
	page, err := strconv.ParseInt(r.FormValue("page"), 10, 64)
	if err != nil {
		page = 1
	}
	hitsPerPage, err := strconv.ParseInt(r.FormValue("hitsPerPage"), 10, 64)
	if err != nil {
		hitsPerPage = coursesPerPage
	}
	sorted, err := strconv.ParseInt(r.FormValue("sort"), 10, 64)
	if err != nil {
		sorted = int64(relevance)
	}
	sortedBy := sortType(sorted)
	semesterInt, err := strconv.ParseInt(r.FormValue("semester"), 10, 64)
	if err != nil {
		semesterInt = int64(teachingBoth)
	}
	semester := TeachingSemester(semesterInt)

	// TODO change language based on URL
	req := Request{
		sessionID:   sessionCookie.Value,
		query:       query,
		indexUID:    courseIndex,
		page:        page,
		hitsPerPage: hitsPerPage,
		lang:        lang,
		sortedBy:    sortedBy,
		semester:    semester,
	}
	return &req, nil
}

func createPageContent(res *Response, req *Request) coursesPage {
	return coursesPage{
		courses:    res.courses,
		page:       int(req.page),
		pageSize:   int(req.hitsPerPage),
		totalPages: int(res.totalPages),
		search:     req.query,
		sortedBy:   req.sortedBy,
		semester:   req.semester,
	}
}

func (s Server) search(req *Request) (*Response, error) {
	// search for courses
	res, err := s.Search.Search(req)
	if err != nil {
		return nil, err
	}
	// retrieve blueprint assignments
	codes := make([]string, len(res.courses))
	for _, course := range res.courses {
		codes = append(codes, course.code)
	}
	assignments, err := s.Data.Blueprint(req.sessionID, codes)
	if err != nil {
		return nil, err
	}
	for i := range res.courses {
		assignment, ok := assignments[res.courses[i].code]
		if ok {
			fmt.Println(assignment)
			res.courses[i].blueprintAssignments = assignment
		}
	}
	return res, nil
}
