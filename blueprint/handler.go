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

func HandleCourseUnassign(w http.ResponseWriter, r *http.Request) {
	err := db.MoveCourse(
		user,
		r.PathValue("code"),
		0, // unassigned year
		0, // semester is not used
		lastPosition, // position is last
	)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO: Render just the necessary for performance
	HandleContent(w, r).Render(r.Context(), w)
}

func HandleCourseAssign(w http.ResponseWriter, r *http.Request) {
	semester, err := strconv.Atoi(r.PathValue("semester"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	year, err := strconv.Atoi(r.FormValue("year"))
	if err != nil {
		http.Error(w, "Unable to parse year", http.StatusBadRequest)
		return
	}
	if semester == int(both) {
		semester, err = strconv.Atoi(r.FormValue("semester"))
		if err != nil {
			http.Error(w, "Unable to parse semester", http.StatusBadRequest)
			return
		}
	}
	err = db.MoveCourse(
		user,
		r.PathValue("code"),
		year,
		semester,
		lastPosition,
	)
	log.Println("Assigning course", r.PathValue("code"), "to year", year, "semester", semester)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO: Render just the necessary for performance
	HandleContent(w, r).Render(r.Context(), w)
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
