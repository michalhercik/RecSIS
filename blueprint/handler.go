package blueprint

import (
	"github.com/a-h/templ"
	"net/http"
	"strconv"
)

const user = 42 // TODO get user from session

func HandleContent(w http.ResponseWriter, r *http.Request) templ.Component {
	data := db.GetData(user)
	return Content(&data)
}

func HandlePage(w http.ResponseWriter, r *http.Request) templ.Component {
	data := db.GetData(user)
	return Page(&data)
}

func HandleUnassignedRemoval(w http.ResponseWriter, r *http.Request) {
	// Remove data from DB
	courseId, _ := strconv.Atoi(r.PathValue("id"))
	db.RemoveUnassigned(user, courseId)
    // Send http response
	w.WriteHeader(http.StatusOK)
}

func HandleLastYearRemoval(w http.ResponseWriter, r *http.Request) {
    year := r.PathValue("year")

    // Update data in DB
    yearInt, _ := strconv.Atoi(year)
    db.RemoveYear(user, yearInt)

    // Send refresh header
    w.Header().Set("HX-Refresh", "true") // htmx will trigger a full page reload
	w.WriteHeader(http.StatusOK)
}

func HandleYearAddition(w http.ResponseWriter, r *http.Request) {
	// Update data in DB
	db.AddYear(user)
	// Send refresh header
    w.Header().Set("HX-Refresh", "true") // htmx will trigger a full page reload
	w.WriteHeader(http.StatusOK)
}
