package blueprint

import (
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
)

const user = 42 // TODO get user from session

func HandleContent(w http.ResponseWriter, r *http.Request) templ.Component {
	data, err := db.BluePrint(user)
	if err != nil {
		log.Println(err)
		return InternalServerErrorContent()
	}
	return Content(data)
}

func HandlePage(w http.ResponseWriter, r *http.Request) {
	data, err := db.BluePrint(user)
	if err != nil {
		log.Println(err)
		InternalServerErrorPage().Render(r.Context(), w)
	} else {
		Page(data).Render(r.Context(), w)
	}
}

func HandleCourseRemoval(w http.ResponseWriter, r *http.Request) {
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
	err = db.RemoveCourse(
		user,
		r.PathValue("code"),
		year,
		semester,
	)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func HandleYearRemoval(w http.ResponseWriter, r *http.Request) {
	if err := db.RemoveYear(user); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	// TODO: Render just credits stats with own sql query for performance
	HandleContent(w, r).
		Render(r.Context(), w)
}

func HandleYearAddition(w http.ResponseWriter, r *http.Request) {
	log.Println("Adding year")
	if err := db.AddYear(user); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO: Render just credits stats with own sql query for performance
	HandleContent(w, r).
		Render(r.Context(), w)
}
