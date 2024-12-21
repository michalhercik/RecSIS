package blueprint

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
)

/**
TODO:
	- handle errors
	- handle logging
	- document functions
*/

const user = 42 // TODO get user from session

type Server struct {
	Data DataManager
}

func (s Server) Register(router *http.ServeMux, prefix string) {
	router.HandleFunc(fmt.Sprintf("GET %s", prefix), s.page)
	router.HandleFunc(fmt.Sprintf("POST %s/year", prefix), s.yearAddition)
	router.HandleFunc(fmt.Sprintf("DELETE %s/year", prefix), s.yearRemoval)
	router.HandleFunc(fmt.Sprintf("POST %s/course/{code}", prefix), s.courseAddition)
	router.HandleFunc(fmt.Sprintf("PATCH %s/course/{id}", prefix), s.courseMovement)
	router.HandleFunc(fmt.Sprintf("PATCH %s/courses", prefix), s.coursesMovement)
	router.HandleFunc(fmt.Sprintf("DELETE %s/course/{id}", prefix), s.courseRemoval)
	router.HandleFunc(fmt.Sprintf("DELETE %s/courses", prefix), s.coursesRemoval)
}

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

func parseSemester(r *http.Request) (SemesterAssignment, error) {
	semesterInt, err := strconv.Atoi(r.FormValue("semester"))
	if err != nil {
		return 0, err
	}
	semester := SemesterAssignment(semesterInt)
	return semester, nil
}

func parseYearSemester(r *http.Request) (int, SemesterAssignment, error) {
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

func parseYearSemesterPosition(r *http.Request) (int, SemesterAssignment, int, error) {
	year, semester, err := parseYearSemester(r)
	if err != nil {
		return year, semester, 0, err
	}
	position, err := strconv.Atoi(r.FormValue("position"))
	return year, semester, position, err
}

func atoiSlice(s []string) ([]int, error) {
	result := make([]int, len(s))
	for i, elem := range s {
		var err error
		result[i], err = strconv.Atoi(elem)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// ===============================================================================================================================
// Page
// ===============================================================================================================================

func (s Server) renderBlueprint(w http.ResponseWriter, r *http.Request) {
	var result templ.Component
	data, err := s.Data.BluePrint(user)
	if err != nil {
		log.Println(err)
		result = InternalServerErrorContent()
	} else {
		result = Content(data)
	}
	result.Render(r.Context(), w)
}

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	var result templ.Component
	data, err := s.Data.BluePrint(user)
	if err != nil {
		log.Println(err)
		result = InternalServerErrorPage()
	} else {
		result = Page(data)
	}
	result.Render(r.Context(), w)
}

// ===============================================================================================================================
// Move Courses
// ===============================================================================================================================

func (s Server) unassignYear(r *http.Request) error {
	year, err := parseYear(r)
	if err != nil {
		return err
	}
	err = s.Data.UnassignYear(user, year)
	return err
}

func (s Server) unassignSemester(r *http.Request) error {
	year, semester, err := parseYearSemester(r)
	if err != nil {
		return err
	}
	err = s.Data.UnassignSemester(user, year, semester)
	return err
}

func (s Server) moveCourses(r *http.Request) error {
	year, semester, position, err := parseYearSemesterPosition(r)
	if err != nil {
		return err
	}
	courses, err := atoiSlice(r.Form["selected"])
	if err != nil {
		return err
	}
	if position == lastPosition {
		err = s.Data.AppendCourses(user, year, semester, courses...)
	} else if position > 0 {
		err = s.Data.InsertCourses(user, year, semester, position, courses...)
	} else {
		err = fmt.Errorf("invalid position %d", position)
	}
	return err
}

func (s Server) coursesMovement(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.FormValue("type") {
	case "year-unassign":
		err = s.unassignYear(r)
	case "semester-unassign":
		err = s.unassignSemester(r)
	case "selected":
		err = s.moveCourses(r)
	default:
		http.Error(w, "Invalid type", http.StatusBadRequest)
	}
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.renderBlueprint(w, r)
}

func (s Server) courseMovement(w http.ResponseWriter, r *http.Request) {
	course, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Unable to parse course ID", http.StatusBadRequest)
		return
	}
	year, semester, position, err := parseYearSemesterPosition(r)
	if err != nil {
		http.Error(w, "Unable to parse parameters", http.StatusBadRequest)
		return
	}
	if position == -1 {
		err = s.Data.AppendCourses(user, year, semester, course)
	} else if position >= 0 {
		err = s.Data.InsertCourses(user, year, semester, position, course)
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
	s.renderBlueprint(w, r)
}

// ===============================================================================================================================
// Remove Courses
// ===============================================================================================================================

func (s Server) removeCoursesByYear(r *http.Request) error {
	year, err := parseYear(r)
	if err != nil {
		return err
	}
	err = s.Data.RemoveCoursesByYear(user, year)
	return err
}

func (s Server) removeCoursesBySemester(r *http.Request) error {
	year, semester, err := parseYearSemester(r)
	if err != nil {
		return err
	}
	err = s.Data.RemoveCoursesBySemester(user, year, semester)
	return err
}

func (s Server) removeCourses(r *http.Request) error {
	r.ParseForm()
	courses, err := atoiSlice(r.Form["selected"])
	if err != nil {
		return err
	}
	log.Println("Removing selected courses", courses)
	err = s.Data.RemoveCourses(user, courses...)
	return err
}

func (s Server) coursesRemoval(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.FormValue("type") {
	case "year":
		err = s.removeCoursesByYear(r)
	case "semester":
		err = s.removeCoursesBySemester(r)
	case "selected":
		err = s.removeCourses(r)
	default:
		http.Error(w, "Invalid type", http.StatusBadRequest)
	}
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.renderBlueprint(w, r)
}

func (s Server) courseRemoval(w http.ResponseWriter, r *http.Request) {
	course, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to parse course ID", http.StatusBadRequest)
		return
	}
	err = s.Data.RemoveCourses(user, course)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO: Render just credits stats with own sql query for performance
	s.renderBlueprint(w, r)
}

// ===============================================================================================================================
// Add Courses
// ===============================================================================================================================

func (s Server) courseAddition(w http.ResponseWriter, r *http.Request) {
	course := r.PathValue("code")
	year, semester, err := parseYearSemester(r)
	if err != nil {
		http.Error(w, "Unable to parse parameters", http.StatusBadRequest)
		return
	}
	courseID, err := s.Data.NewCourse(user, course, year, semester)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": courseID})
}

// ===============================================================================================================================
// Remove Year
// ===============================================================================================================================

func (s Server) yearRemoval(w http.ResponseWriter, r *http.Request) {
	if err := s.Data.RemoveYear(user); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	// TODO: Render just credits stats with own sql query for performance
	s.renderBlueprint(w, r)
}

// ===============================================================================================================================
// Add Year
// ===============================================================================================================================

func (s Server) yearAddition(w http.ResponseWriter, r *http.Request) {
	log.Println("Adding year")
	if err := s.Data.AddYear(user); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO: Render just credits stats with own sql query for performance
	s.renderBlueprint(w, r)
}
