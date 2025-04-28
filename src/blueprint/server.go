package blueprint

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/language"
)

/**
TODO:
	- handle errors
	- handle logging
	- document functions
*/

type Server struct {
	router *http.ServeMux
	Data   DataManager
	Auth   Authentication
}

func (s *Server) Init() {
	s.router = http.NewServeMux()
	s.router.HandleFunc("GET /", s.page)
	s.router.HandleFunc("POST /year", s.yearAddition)
	s.router.HandleFunc("DELETE /year", s.yearRemoval)
	s.router.HandleFunc("PATCH /{year}/{semester}", s.foldSemester)
	s.router.HandleFunc("POST /course/{code}", s.courseAddition)
	s.router.HandleFunc("PATCH /course/{id}", s.courseMovement)
	s.router.HandleFunc("PATCH /courses", s.coursesMovement)
	s.router.HandleFunc("DELETE /course/{id}", s.courseRemoval)
	s.router.HandleFunc("DELETE /courses", s.coursesRemoval)
}

func (s Server) Router() http.Handler {
	return s.router
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

func parseLanguage(r *http.Request) (text, error) {
	lang := language.FromContext(r.Context())
	return texts[lang], nil
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

func (s Server) renderBlueprint(w http.ResponseWriter, r *http.Request, t text) {
	var result templ.Component
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	data, err := s.Data.Blueprint(userID, DBLang(t.Language))
	if err != nil {
		log.Println(err)
		result = InternalServerErrorContent(t)
	} else {
		result = Content(data, t)
	}
	result.Render(r.Context(), w)
}

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	var result templ.Component
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	data, err := s.Data.Blueprint(userID, DBLang(t.Language))
	if err != nil {
		log.Println(err)
		result = InternalServerErrorPage(t)
	} else {
		result = Page(data, t)
	}
	result.Render(r.Context(), w)
}

// ===============================================================================================================================
// Move Courses
// ===============================================================================================================================

func (s Server) unassignYear(r *http.Request, userSession string) error {
	year, err := parseYear(r)
	if err != nil {
		return err
	}
	err = s.Data.UnassignYear(userSession, year)
	return err
}

func (s Server) unassignSemester(r *http.Request, userSession string) error {
	year, semester, err := parseYearSemester(r)
	if err != nil {
		return err
	}
	err = s.Data.UnassignSemester(userSession, year, semester)
	return err
}

func (s Server) moveCourses(r *http.Request, userSession string) error {
	year, semester, position, err := parseYearSemesterPosition(r)
	if err != nil {
		return err
	}
	courses, err := atoiSlice(r.Form["selected"])
	if err != nil {
		return err
	}
	if position == lastPosition {
		err = s.Data.AppendCourses(userSession, year, semester, courses...)
	} else if position > 0 {
		err = s.Data.InsertCourses(userSession, year, semester, position, courses...)
	} else {
		err = fmt.Errorf("invalid position %d", position)
	}
	return err
}

func (s Server) coursesMovement(w http.ResponseWriter, r *http.Request) {
	var err error
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	switch r.FormValue("type") {
	case yearUnassign:
		err = s.unassignYear(r, userID)
	case semesterUnassign:
		err = s.unassignSemester(r, userID)
	case selectedMove:
		err = s.moveCourses(r, userID)
	default:
		http.Error(w, "Invalid type", http.StatusBadRequest)
	}
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	text, err := parseLanguage(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.renderBlueprint(w, r, text)
}

func (s Server) courseMovement(w http.ResponseWriter, r *http.Request) {
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
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
		err = s.Data.AppendCourses(userID, year, semester, course)
	} else if position >= 0 {
		err = s.Data.InsertCourses(userID, year, semester, position, course)
	} else {
		http.Error(w, "Invalid position", http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	text, err := parseLanguage(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO: Render just the necessary for performance
	s.renderBlueprint(w, r, text)
}

// ===============================================================================================================================
// Remove Courses
// ===============================================================================================================================

func (s Server) removeCoursesByYear(r *http.Request, session string) error {
	year, err := parseYear(r)
	if err != nil {
		return err
	}
	err = s.Data.RemoveCoursesByYear(session, year)
	return err
}

func (s Server) removeCoursesBySemester(r *http.Request, session string) error {
	year, semester, err := parseYearSemester(r)
	if err != nil {
		return err
	}
	err = s.Data.RemoveCoursesBySemester(session, year, semester)
	return err
}

func (s Server) removeCourses(r *http.Request, session string) error {
	r.ParseForm()
	courses, err := atoiSlice(r.Form["selected"])
	if err != nil {
		return err
	}
	err = s.Data.RemoveCourses(session, courses...)
	return err
}

func (s Server) coursesRemoval(w http.ResponseWriter, r *http.Request) {
	var err error
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	switch r.FormValue("type") {
	case yearRemove:
		err = s.removeCoursesByYear(r, userID)
	case semesterRemove:
		err = s.removeCoursesBySemester(r, userID)
	case selectedRemove:
		err = s.removeCourses(r, userID)
	default:
		http.Error(w, "Invalid type", http.StatusBadRequest)
	}
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	text, err := parseLanguage(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.renderBlueprint(w, r, text)
}

func (s Server) courseRemoval(w http.ResponseWriter, r *http.Request) {
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	course, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to parse course ID", http.StatusBadRequest)
		return
	}
	err = s.Data.RemoveCourses(userID, course)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	text, err := parseLanguage(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO: Render just credits stats with own sql query for performance
	s.renderBlueprint(w, r, text)
}

// ===============================================================================================================================
// Add Courses
// ===============================================================================================================================

type courseAdditionPresenter func(insertedCourseInfo) templ.Component

func (s Server) courseAddition(w http.ResponseWriter, r *http.Request) {
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	course := r.PathValue("code")
	year, semester, err := parseYearSemester(r)
	if err != nil {
		log.Printf("courseAddition error: %v", err)
		return
	}
	reqSourceInt, err := strconv.Atoi(r.FormValue("requestSource"))
	if err != nil {
		reqSourceInt = int(sourceNone)
		log.Printf("courseAddition warning: %v", err)
	}
	reqSource := courseAdditionRequestSource(reqSourceInt)
	var presenter courseAdditionPresenter = DefaultCourseAdditionPresenter
	switch reqSource {
	case sourceBlueprint:
		// TODO: implement
		//
		// Example:
		//	file: blueprint/blueprint.templ
		//		templ Ribbon(insertInfo insertedCourseInfo) {
		//			<div class="ribbon"> {{insertInfo.courseID}} </div>
		//  	}
		//
		//	file: blueprint/blueprint.go
		//		presenter = Ribbon
	case sourceCourseDetail:
		// TODO: implement
	case sourceDegreePlan:
		// TODO: implement
	}
	// TODO: use this
	// text, err := parseLanguage(r)
	// if err != nil {
	// 	log.Println(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	courseID, err := s.Data.NewCourse(userID, course, year, semester)
	if err != nil {
		log.Println(err)
		return
	}
	insertInfo := insertedCourseInfo{courseID: courseID, academicYear: year, semester: semester}
	presenter(insertInfo).Render(r.Context(), w)
}

// ===============================================================================================================================
// Remove Year
// ===============================================================================================================================

func (s Server) yearRemoval(w http.ResponseWriter, r *http.Request) {
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
	if err := s.Data.RemoveYear(userID); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	text, err := parseLanguage(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO: Render just credits stats with own sql query for performance
	s.renderBlueprint(w, r, text)
}

// ===============================================================================================================================
// Add Year
// ===============================================================================================================================

func (s Server) yearAddition(w http.ResponseWriter, r *http.Request) {
	// sessionCookie, err := r.Cookie("recsis_session_key")
	// if err != nil {
	// 	http.Error(w, "unknown student", http.StatusBadRequest)
	// 	log.Printf("HandlePage error: %v", err)
	// 	return
	// }
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if err := s.Data.AddYear(userID); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	text, err := parseLanguage(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO: Render just credits stats with own sql query for performance
	s.renderBlueprint(w, r, text)
}

func (s Server) foldSemester(w http.ResponseWriter, r *http.Request) {
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	year, err := strconv.Atoi(r.PathValue("year"))
	if err != nil {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		log.Printf("foldSemester error: %v", err)
		return
	}
	semester, err := SemesterAssignmentFromString(r.PathValue("semester"))
	if err != nil {
		http.Error(w, "Invalid semester", http.StatusBadRequest)
		log.Printf("foldSemester error: %v", err)
		return
	}
	folded, err := strconv.ParseBool(r.FormValue("folded"))
	if err != nil {
		http.Error(w, "Invalid folded value", http.StatusBadRequest)
		log.Printf("foldSemester error: %v", err)
		return
	}
	fmt.Println(userID, year, semester, folded)
	err = s.Data.FoldSemester(userID, year, semester, folded)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// text, err := parseLanguage(r)
	// if err != nil {
	// 	log.Println(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// s.renderBlueprint(w, r, text)
}
