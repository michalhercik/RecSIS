package blueprint

import (
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	
    "github.com/michalhercik/RecSIS/courses"
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
	course, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to parse course ID", http.StatusBadRequest)
		return
	}
	err = db.RemoveCourse(user, course)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO: Render just credits stats with own sql query for performance
	HandleContent(w, r).
		Render(r.Context(), w)
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

func HandleCourseAddition(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	origin := r.PostFormValue("origin")
	course := r.PathValue("code")
	year, err := strconv.Atoi(r.FormValue("year"))
	if err != nil {
		http.Error(w, "Unable to parse year", http.StatusBadRequest)
		return
	}
	semesterInt, err := strconv.Atoi(r.FormValue("semester"))
	// TODO: check validity of semesterInt
	semester := SemesterPosition(semesterInt)
	if err != nil {
		http.Error(w, "Unable to parse semester", http.StatusBadRequest)
		return
	}
	err = db.InsertCourse(
		user,
		course,
		year,
		semester,
	)
	if err != nil {
		http.Error(w, "Unable to parse semester", http.StatusBadRequest)
	}
	switch origin {
	case "courses":
		courses.BlueprintAssignment(nil, course).Render(r.Context(), w)
		//http.Error(w, "Not implemented", http.StatusNotImplemented)
	case "coursedetail":
		http.Error(w, "Not implemented", http.StatusNotImplemented)
	case "degreeplan":
		http.Error(w, "Not implemented", http.StatusNotImplemented)
	default:
		http.Error(w, "Not implemented", http.StatusNotImplemented)
	}
}

func HandleCourseMovement(w http.ResponseWriter, r *http.Request) {
	course, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Unable to parse course ID", http.StatusBadRequest)
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
	semesterInt, err := strconv.Atoi(r.FormValue("semester"))
	// TODO: check validity of semesterInt
	semester := SemesterPosition(semesterInt)
	if err != nil {
		http.Error(w, "Unable to parse semester", http.StatusBadRequest)
		return
	}
	err = db.AppendCourse(
		user,
		course,
		year,
		semester,
	)
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
