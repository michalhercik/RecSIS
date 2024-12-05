package courses

import (
	"github.com/a-h/templ"
	"net/http"
	"strconv"
)

func getDefaultQuery() query {
	return query{
		user:       user,
		startIndex: 0,
		maxCount:   coursesPerPage,
		search:     "",
		sorted:     recommended,
	}
}

func HandleContent(w http.ResponseWriter, r *http.Request) templ.Component {
	recommendedCourses, _ := db.Courses(getDefaultQuery())
	return Content(&recommendedCourses)
}

func HandlePage(w http.ResponseWriter, r *http.Request) templ.Component {
	recommendedCourses, _ := db.Courses(getDefaultQuery())
	return Page(&recommendedCourses)
}

func HandlePaging(w http.ResponseWriter, r *http.Request) {
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
    coursesPage, _ := db.Courses(query)

    // Render search results
    Courses(&coursesPage).Render(r.Context(), w)
}

func sortTypeFromString(st string) sortType {
    switch st {
    case "relevance":
        return relevance
    case "recommended":
        return recommended
    case "rating":
        return rating
    case "most_popular":
        return mostPopular
    case "newest":
        return newest
    default:
        return relevance
    }
}

func HandleSearch(w http.ResponseWriter, r *http.Request) {
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
    coursesPage, _ := db.Courses(query)

    // Render search results
    Courses(&coursesPage).Render(r.Context(), w)
}
