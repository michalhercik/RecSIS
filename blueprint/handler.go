package blueprint

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
)

// ===============================================================================================================================
// Utils
// ===============================================================================================================================

func parseYear(r *http.Request) (int, error) {
	year, err := strconv.Atoi(r.FormValue("year"))
	if err != nil {
		return year, err
	}
	return year, nil
}

func parseSemester(r *http.Request) (SemesterPosition, error) {
	semesterInt, err := strconv.Atoi(r.FormValue("semester"))
	if err != nil {
		return 0, err
	}
	semester := SemesterPosition(semesterInt)
	return semester, nil
}

func parseYearSemester(r *http.Request) (int, SemesterPosition, error) {
	year, err := parseYear(r)
	if err != nil {
		return 0, 0, err
	}
	semester, err := parseSemester(r)
	if err != nil {
		return year, 0, err
	}
	return year, semester, nil
}

// ===============================================================================================================================
// Page
// ===============================================================================================================================

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

// ===============================================================================================================================
// Move Courses
// ===============================================================================================================================

func unassignYear(r *http.Request) error {
	year, err := parseYear(r)
	if err != nil {
		return err
	}
	err = db.UnassignYear(user, year)
	return err
}

func unassignSemester(r *http.Request) error {
	year, semester, err := parseYearSemester(r)
	if err != nil {
		return err
	}
	err = db.UnassignSemester(user, year, semester)
	return err
}

// TODO: Implement, connect to db, add more cases
func HandleCoursesMovement(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.FormValue("type") {
	case "year-unassign":
		err = unassignYear(r)
	case "semester-unassign":
		err = unassignSemester(r)
	case "selected":
		log.Println("Unassigning selected courses")
	default:
		http.Error(w, "Invalid type", http.StatusBadRequest)
	}

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	HandleContent(w, r).
		Render(r.Context(), w)
}

// ===============================================================================================================================
// Remove Courses
// ===============================================================================================================================

func removeCoursesByYear(r *http.Request) error {
	year, err := parseYear(r)
	if err != nil {
		return err
	}
	err = db.RemoveCoursesByYear(user, year)
	return err
}

func removeCoursesBySemester(r *http.Request) error {
	year, semester, err := parseYearSemester(r)
	if err != nil {
		return err
	}
	err = db.RemoveCoursesBySemester(user, year, semester)
	return err
}

// TODO: Implement, connect to db, add more cases
func HandleCoursesRemoval(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.FormValue("type") {
	case "year":
		err = removeCoursesByYear(r)
	case "semester":
		err = removeCoursesBySemester(r)
	case "selected":
		log.Println("Removing selected courses")
	default:
		http.Error(w, "Invalid type", http.StatusBadRequest)
	}
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
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
	err = db.RemoveCourses(user, course)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO: Render just credits stats with own sql query for performance
	HandleContent(w, r).
		Render(r.Context(), w)
}

// ===============================================================================================================================
// ...
// ===============================================================================================================================

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
	course := r.PathValue("code")
	year, err := strconv.Atoi(r.PostFormValue("year"))
	if err != nil {
		http.Error(w, "Unable to parse year", http.StatusBadRequest)
		return
	}
	semesterInt, err := strconv.Atoi(r.PostFormValue("semester"))
	// TODO: check validity of semesterInt
	semester := SemesterPosition(semesterInt)
	if err != nil {
		http.Error(w, "Unable to parse semester", http.StatusBadRequest)
		return
	}

	courseID, err := db.InsertCourse(
		user,
		course,
		year,
		semester,
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": courseID})
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
		err = db.AppendCourses(user, year, semester, course)
	} else if position >= 0 {
		err = db.MoveCourses(user, year, semester, position, course)
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
