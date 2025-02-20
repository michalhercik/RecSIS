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
	router.HandleFunc(fmt.Sprintf("GET %s", prefix), s.page)
	// router.HandleFunc(fmt.Sprintf("GET %s/page", prefix), s.paging)
	router.HandleFunc(fmt.Sprintf("GET %s/search", prefix), s.content)
	router.HandleFunc(fmt.Sprintf("GET %s/quicksearch", prefix), s.quickSearch)
}

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	req := parseQueryRequest(r)
	res, err := s.search(req)
	if err != nil {
		log.Printf("quickSearch: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	coursesPage := createPageContent(res, req)
	Page(&coursesPage).Render(r.Context(), w)
}

func (s Server) content(w http.ResponseWriter, r *http.Request) {
	req := parseQueryRequest(r)
	res, err := s.search(req)
	if err != nil {
		log.Printf("quickSearch: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	coursesPage := createPageContent(res, req)
	Courses(&coursesPage).Render(r.Context(), w)
}

func (s Server) quickSearch(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("search")
	req := QuickRequest{
		query:    query,
		indexUID: courseIndex,
		limit:    5,
		offset:   0,
		lang:     czech,
	}
	res, err := s.Search.QuickSearch(&req)
	if err != nil {
		log.Printf("quickSearch: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = res
	//TODO: render res
}

func parseQueryRequest(r *http.Request) Request {
	query := r.FormValue("search")
	page, err := strconv.ParseInt(r.FormValue("page"), 10, 64)
	if err != nil {
		page = 1
	}
	hitsPerPage, err := strconv.ParseInt(r.FormValue("hitsPerPage"), 10, 64)
	if err != nil {
		hitsPerPage = coursesPerPage
	}
	sorted, err := strconv.ParseInt(r.FormValue("sort"), 10, 32)
	if err != nil {
		sorted = 0
	}
	sortedBy := sortType(sorted)

	req := Request{
		query:       query,
		indexUID:    courseIndex,
		page:        page,
		hitsPerPage: hitsPerPage,
		lang:        czech,
		sortedBy:    sortedBy,
	}
	return req
}

func createPageContent(res *Response, req Request) coursesPage {
	return coursesPage{
		courses:    res.courses,
		page:       int(req.page),
		pageSize:   int(req.hitsPerPage),
		totalPages: int(res.totalPages),
		search:     req.query,
		sortedBy:   req.sortedBy,
	}
}

func (s Server) search(req Request) (*Response, error) {
	// search for courses
	res, err := s.Search.Search(&req)
	if err != nil {
		return nil, err
	}
	// retrieve blueprint assignments
	codes := make([]string, len(res.courses))
	for _, course := range res.courses {
		codes = append(codes, course.code)
	}
	assignments, err := s.Data.Blueprint(user, codes)
	if err != nil {
		return nil, err
	}
	for i := range res.courses {
		assignment, ok := assignments[res.courses[i].code]
		if ok {
			res.courses[i].blueprintAssignments = assignment
		}
	}
	return res, nil
}

// func getDefaultQuery() query {
// 	return query{
// 		user:       user,
// 		startIndex: 0,
// 		maxCount:   coursesPerPage,
// 		search:     "",
// 		sorted:     recommended,
// 	}
// }

// func (s Server) page(w http.ResponseWriter, r *http.Request) {
// 	recommendedCourses, err := s.Data.Courses(getDefaultQuery())
// 	if err != nil {
// 		log.Printf("HandlePage: %v", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 	} else {
// 		Page(&recommendedCourses).Render(r.Context(), w)
// 	}
// }

// func (s Server) paging(w http.ResponseWriter, r *http.Request) {
// 	// Parse query parameters from URL
// 	user, _ := strconv.Atoi(r.URL.Query().Get("user"))
// 	startIndex, _ := strconv.Atoi(r.URL.Query().Get("startIndex"))
// 	maxCount, _ := strconv.Atoi(r.URL.Query().Get("maxCount"))
// 	search := r.URL.Query().Get("search")
// 	sorted, _ := strconv.Atoi(r.URL.Query().Get("sorted"))

// 	// Create query from input
// 	query := query{
// 		user:       user,
// 		startIndex: startIndex,
// 		maxCount:   maxCount,
// 		search:     search,
// 		sorted:     sortType(sorted),
// 	}

// 	// Get result from search
// 	coursesPage, _ := s.Data.Courses(query)

// 	// Render search results
// 	Courses(&coursesPage).Render(r.Context(), w)
// }

// func sortTypeFromString(st string) sortType {
// 	switch st {
// 	case op_relevance:
// 		return relevance
// 	case op_recommended:
// 		return recommended
// 	case op_rating:
// 		return rating
// 	case op_mostPopular:
// 		return mostPopular
// 	case op_newest:
// 		return newest
// 	default:
// 		return relevance
// 	}
// }

// func (s Server) search(w http.ResponseWriter, r *http.Request) {
// 	// Parse query parameters from URL
// 	user, _ := strconv.Atoi(r.URL.Query().Get("user"))
// 	search := r.URL.Query().Get("search")
// 	sorted := sortTypeFromString(r.URL.Query().Get("sort"))

// 	// Create query from input
// 	query := query{
// 		user:       user,
// 		startIndex: 0,
// 		maxCount:   coursesPerPage,
// 		search:     search,
// 		sorted:     sortType(sorted),
// 	}

// 	// Get result from search
// 	coursesPage, _ := s.Data.Courses(query)

// 	// Render search results
// 	Courses(&coursesPage).Render(r.Context(), w)
// }
