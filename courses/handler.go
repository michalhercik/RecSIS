package courses

import (
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

func HandlePage(w http.ResponseWriter, r *http.Request) {
	recommendedCourses, _ := db.Courses(getDefaultQuery())
	Page(&recommendedCourses).Render(r.Context(), w)
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

func HandleBlueprintAddition(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters from URL
	code := r.PathValue("code")

	// Make data changes
	assignments, err := db.AddCourseToBlueprint(user, code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Render the result
	BlueprintAssignment(assignments, code).Render(r.Context(), w)
}

func HandleBlueprintRemoval(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters from URL
	code := r.PathValue("code")

	// Make data changes
	err := db.RemoveCourseFromBlueprint(user, code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Render the result
	BlueprintAssignment([]Assignment{}, code).Render(r.Context(), w)
}
