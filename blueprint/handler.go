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

// TODO: Implement, connect to db, add more cases
func HandleCoursesMovement(w http.ResponseWriter, r *http.Request) {
	callType := r.FormValue("type")
	year, err := strconv.Atoi(r.FormValue("year"))
	if err != nil {
		http.Error(w, "Unable to parse year", http.StatusBadRequest)
		return
	}
	switch callType {
	case "year-unassign":
		log.Println("Unassigning all courses from year", year)
	case "semester-unassign":
		semesterInt, err := strconv.Atoi(r.FormValue("semester"))
		if err != nil {
			http.Error(w, "Unable to parse semester", http.StatusBadRequest)
			return
		}
		log.Println("Unassigning all courses from year/semester", year, semesterInt)
	case "selected":
		semesterInt, err := strconv.Atoi(r.FormValue("semester"))
		if err != nil {
			http.Error(w, "Unable to parse semester", http.StatusBadRequest)
			return
		}
		r.ParseForm()
		courses := r.Form["selected"]
		log.Println("Moving selected courses to year/semester", courses, year, semesterInt)
	default:
		 http.Error(w, "Invalid type", http.StatusBadRequest)
	}

	HandleContent(w, r).
		Render(r.Context(), w)
}

// TODO: Implement, connect to db, add more cases
func HandleCoursesRemoval(w http.ResponseWriter, r *http.Request) {
	callType := r.FormValue("type")
	switch callType {
	case "year":
		year, err := strconv.Atoi(r.FormValue("year"))
		if err != nil {
			http.Error(w, "Unable to parse year", http.StatusBadRequest)
			return
		}
		log.Println("Removing all courses from year", year)
	case "semester":
		year, err := strconv.Atoi(r.FormValue("year"))
		if err != nil {
			http.Error(w, "Unable to parse year", http.StatusBadRequest)
			return
		}
		semesterInt, err := strconv.Atoi(r.FormValue("semester"))
		if err != nil {
			http.Error(w, "Unable to parse semester", http.StatusBadRequest)
			return
		}
		log.Println("Removing all courses from semester", year, semesterInt)
	case "selected":
		r.ParseForm()
		courses := r.Form["selected"]
		log.Println("Removing selected courses", courses)
	default:
		 http.Error(w, "Invalid type", http.StatusBadRequest)
	}

	HandleContent(w, r).
		Render(r.Context(), w)
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
	log.Println("Adding course from", origin)
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
		log.Println(r.FormValue("year"))
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
	position, err := strconv.Atoi(r.FormValue("position"))
	if err != nil {
		http.Error(w, "Unable to parse position", http.StatusBadRequest)
		return
	}
	if position == -1 {
		err = db.AppendCourse(user, course, year, semester)
	} else if position >= 0 {
		err = db.MoveCourse(user, course, year, semester, position)
	} else {
		http.Error(w, "Invalid position", http.StatusBadRequest)
		return
	}
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
