package courses

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	Data DataManager
}

func (s Server) Register(router *http.ServeMux, prefix string) {
	router.HandleFunc(fmt.Sprintf("GET %s", prefix), s.page)
	router.HandleFunc(fmt.Sprintf("GET %s/page", prefix), s.paging)
	router.HandleFunc(fmt.Sprintf("GET %s/search", prefix), s.search)
}

func getDefaultQuery() query {
	return query{
		user:       user,
		startIndex: 0,
		maxCount:   coursesPerPage,
		search:     "",
		sorted:     recommended,
	}
}

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	recommendedCourses, err := s.Data.Courses(getDefaultQuery())
	if err != nil {
		log.Printf("HandlePage: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		Page(&recommendedCourses).Render(r.Context(), w)
	}
}

func (s Server) paging(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters from URL
	user, _ := strconv.Atoi(r.URL.Query().Get("user"))
	startIndex, _ := strconv.Atoi(r.URL.Query().Get("startIndex"))
	maxCount, _ := strconv.Atoi(r.URL.Query().Get("maxCount"))
	search := r.URL.Query().Get("search")
	sorted, _ := strconv.Atoi(r.URL.Query().Get("sorted"))

	// Create query from input
	query := query{
		user:       user,
		startIndex: startIndex,
		maxCount:   maxCount,
		search:     search,
		sorted:     sortType(sorted),
	}

	// Get result from search
	coursesPage, _ := s.Data.Courses(query)

	// Render search results
	Courses(&coursesPage).Render(r.Context(), w)
}

func sortTypeFromString(st string) sortType {
	switch st {
	case op_relevance:
		return relevance
	case op_recommended:
		return recommended
	case op_rating:
		return rating
	case op_mostPopular:
		return mostPopular
	case op_newest:
		return newest
	default:
		return relevance
	}
}

func (s Server) search(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters from URL
	user, _ := strconv.Atoi(r.URL.Query().Get("user"))
	search := r.URL.Query().Get("search")
	sorted := sortTypeFromString(r.URL.Query().Get("sort"))

	// Create query from input
	query := query{
		user:       user,
		startIndex: 0,
		maxCount:   coursesPerPage,
		search:     search,
		sorted:     sortType(sorted),
	}

	// Get result from search
	coursesPage, _ := s.Data.Courses(query)

	// Render search results
	Courses(&coursesPage).Render(r.Context(), w)
}
