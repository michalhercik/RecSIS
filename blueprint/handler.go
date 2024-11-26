package blueprint

import (
	"github.com/a-h/templ"
    "github.com/michalhercik/RecSIS/database"
	"net/http"
	"strconv"
)

func HandleContent(w http.ResponseWriter, r *http.Request) templ.Component {
	data := database.GetBlueprintData(database.User)
	return Content(&data)
}

func HandlePage(w http.ResponseWriter, r *http.Request) templ.Component {
	data := database.GetBlueprintData(database.User)
	return Page(&data)
}

func HandleUnassignedRemoval(w http.ResponseWriter, r *http.Request) {
	// Remove data from DB
	courseId, _ := strconv.Atoi(r.PathValue("id"))
	database.BlueprintRemoveUnassigned(database.User, courseId)
    // Send http response
	w.WriteHeader(http.StatusOK)
}

func HandleLastYearRemoval(w http.ResponseWriter, r *http.Request) {
    year := r.PathValue("year")

    // Remove data from DB
    yearInt, _ := strconv.Atoi(year)
    database.BlueprintRemoveYear(database.User, yearInt)

    // Send refresh header
    w.Header().Set("HX-Refresh", "true") // htmx will trigger a full page reload
	w.WriteHeader(http.StatusOK)
}

func HandleYearAddition(w http.ResponseWriter, r *http.Request) {
	// Update DB
	database.BlueprintAddYear(database.User)
	// Send refresh header
    w.Header().Set("HX-Refresh", "true") // htmx will trigger a full page reload
	w.WriteHeader(http.StatusOK)
}
