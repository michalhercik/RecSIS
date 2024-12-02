package blueprint

import (
	"net/http"
	"strconv"

	"github.com/a-h/templ"
)

const user = 42 // TODO get user from session

func HandleContent(w http.ResponseWriter, r *http.Request) templ.Component {
	data, err := db.BluePrint(user)
	if err != nil {
		return InternalServerErrorContent()
	}
	return Content(data)
}

func HandlePage(w http.ResponseWriter, r *http.Request) templ.Component {
	data, err := db.BluePrint(user)
	if err != nil {
		return InternalServerErrorPage()
	}
	return Page(data)
}

func HandleCourseRemoval(w http.ResponseWriter, r *http.Request) {
	// Remove data from DB
	year, err := strconv.Atoi(r.PathValue("year"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	semester, err := strconv.Atoi(r.PathValue("semester"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	db.RemoveCourse(
		user,
		r.PathValue("code"),
		year,
		semester,
	)
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
