package blueprint

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/language"
)

/**
TODO:
	- handle errors
	- handle logging
	- document functions
*/

//================================================================================
// Server Type
//================================================================================

type Server struct {
	Auth   Authentication
	Data   DBManager
	Page   Page
	router *http.ServeMux
}

func (s Server) Router() http.Handler {
	return s.router
}

type Authentication interface {
	UserID(r *http.Request) (string, error)
}

type Page interface {
	View(main templ.Component, lang language.Language, title string, userID string) templ.Component
}

//================================================================================
// Routing
//================================================================================

func (s *Server) Init() {
	s.router = http.NewServeMux()
	s.router.HandleFunc("GET /", s.page)
	s.router.HandleFunc("PATCH /course/{id}", s.courseMovement)
	s.router.HandleFunc("PATCH /courses", s.coursesMovement)
	s.router.HandleFunc("DELETE /course/{id}", s.courseRemoval)
	s.router.HandleFunc("DELETE /courses", s.coursesRemoval)
	s.router.HandleFunc("POST /year", s.yearAddition)
	s.router.HandleFunc("DELETE /year", s.yearRemoval)
	s.router.HandleFunc("PATCH /fold", s.foldSemester)
}

//================================================================================
// Page
//================================================================================

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	var result templ.Component
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	data, err := s.Data.blueprint(userID, t.language)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		generateWarnings(data, t)
		result = Content(data, t)
	}
	s.Page.View(result, lang, t.pageTitle, userID).Render(r.Context(), w)
}

func (s Server) renderBlueprint(w http.ResponseWriter, r *http.Request, t text) {
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	data, err := s.Data.blueprint(userID, t.language)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	generateWarnings(data, t)
	Content(data, t).Render(r.Context(), w)
}

//================================================================================
// Move Courses
//================================================================================

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
		err = s.Data.appendCourses(userID, year, semester, course)
	} else if position >= 0 {
		err = s.Data.insertCourses(userID, year, semester, position, course)
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

func (s Server) coursesMovement(w http.ResponseWriter, r *http.Request) {
	var err error
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	switch r.FormValue(typeParam) {
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

func (s Server) unassignYear(r *http.Request, userSession string) error {
	year, err := parseYear(r)
	if err != nil {
		return err
	}
	err = s.Data.unassignYear(userSession, year)
	return err
}

func (s Server) unassignSemester(r *http.Request, userSession string) error {
	year, semester, err := parseYearSemester(r)
	if err != nil {
		return err
	}
	err = s.Data.unassignSemester(userSession, year, semester)
	return err
}

func (s Server) moveCourses(r *http.Request, userSession string) error {
	year, semester, position, err := parseYearSemesterPosition(r)
	if err != nil {
		return err
	}
	courses, err := atoiSlice(r.Form[checkboxName])
	if err != nil {
		return err
	}
	if position == lastPosition {
		err = s.Data.appendCourses(userSession, year, semester, courses...)
	} else if position > 0 {
		err = s.Data.insertCourses(userSession, year, semester, position, courses...)
	} else {
		err = fmt.Errorf("invalid position %d", position)
	}
	return err
}

// ================================================================================
// Remove Courses
// ================================================================================

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
	err = s.Data.removeCourses(userID, course)
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

func (s Server) coursesRemoval(w http.ResponseWriter, r *http.Request) {
	var err error
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	switch r.FormValue(typeParam) {
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
func (s Server) removeCoursesByYear(r *http.Request, session string) error {
	year, err := parseYear(r)
	if err != nil {
		return err
	}
	err = s.Data.removeCoursesByYear(session, year)
	return err
}

func (s Server) removeCoursesBySemester(r *http.Request, session string) error {
	year, semester, err := parseYearSemester(r)
	if err != nil {
		return err
	}
	err = s.Data.removeCoursesBySemester(session, year, semester)
	return err
}

func (s Server) removeCourses(r *http.Request, session string) error {
	r.ParseForm()
	courses, err := atoiSlice(r.Form[checkboxName])
	if err != nil {
		return err
	}
	err = s.Data.removeCourses(session, courses...)
	return err
}

//================================================================================
// Add/Remove Year
//================================================================================

func (s Server) yearAddition(w http.ResponseWriter, r *http.Request) {
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if err := s.Data.addYear(userID); err != nil {
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

func (s Server) yearRemoval(w http.ResponseWriter, r *http.Request) {
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
	unassign := r.FormValue(unassignParam)
	shouldUnassign, err := strconv.ParseBool(unassign)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	year, err := parseYear(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := s.Data.removeYear(userID, year, shouldUnassign); err != nil {
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

//================================================================================
// Fold Semester
//================================================================================

func (s Server) foldSemester(w http.ResponseWriter, r *http.Request) {
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	year, err := parseYear(r)
	if err != nil {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		log.Printf("foldSemester error: %v", err)
		return
	}
	semester, err := parseSemester(r)
	if err != nil {
		http.Error(w, "Invalid semester", http.StatusBadRequest)
		log.Printf("foldSemester error: %v", err)
		return
	}
	folded, err := strconv.ParseBool(r.FormValue(foldedParam))
	if err != nil {
		http.Error(w, "Invalid folded value", http.StatusBadRequest)
		log.Printf("foldSemester error: %v", err)
		return
	}
	err = s.Data.foldSemester(userID, year, semester, folded)
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
	// TODO: Render just targeted table for performance
	s.renderBlueprint(w, r, text)
}

//================================================================================
// Utils
//================================================================================

func parseYear(r *http.Request) (int, error) {
	year, err := strconv.Atoi(r.FormValue(yearParam))
	if err != nil {
		return year, err
	}
	return year, nil
}

func parseSemester(r *http.Request) (semesterAssignment, error) {
	semesterInt, err := strconv.Atoi(r.FormValue(semesterParam))
	if err != nil {
		return 0, err
	}
	semester := semesterAssignment(semesterInt)
	return semester, nil
}

func parseYearSemester(r *http.Request) (int, semesterAssignment, error) {
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

func parseYearSemesterPosition(r *http.Request) (int, semesterAssignment, int, error) {
	year, semester, err := parseYearSemester(r)
	if err != nil {
		return year, semester, 0, err
	}
	position, err := strconv.Atoi(r.FormValue(positionParam))
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

func generateWarnings(bp *blueprintPage, t text) {
	// check if the course is assigned in a correct semester
	for _, year := range bp.years {
		winter := year.winter
		if !winter.folded {
			courses := winter.courses
			for ci := range courses {
				if courses[ci].semester == teachingSummerOnly {
					courses[ci].warnings = append(courses[ci].warnings, t.wWrongAssignWinter)
				}
			}
		}
		summer := year.summer
		if !summer.folded {
			courses := summer.courses
			for ci := range courses {
				if courses[ci].semester == teachingWinterOnly {
					courses[ci].warnings = append(courses[ci].warnings, t.wWrongAssignSummer)
				}
			}
		}
	}
	// check if the course is assigned more than once
	// TODO: ignore the course if it can be completed more than once
	for _, year1 := range bp.years {
		winter1 := year1.winter.courses
		for ci1 := range winter1 {
			duplicates := make([]struct {
				year     int
				semester string
			}, 0)
			for y2, year2 := range bp.years {
				winter2 := year2.winter.courses
				for ci2 := range winter2 {
					if winter1[ci1].code == winter2[ci2].code {
						duplicates = append(duplicates, struct {
							year     int
							semester string
						}{year: y2, semester: t.winter})
					}
				}
				summer2 := year2.summer.courses
				for ci2 := range summer2 {
					if winter1[ci1].code == summer2[ci2].code {
						duplicates = append(duplicates, struct {
							year     int
							semester string
						}{year: y2, semester: t.summer})
					}
				}
			}
			if len(duplicates) > 1 {
				warning := t.wAssignedMoreThanOnce + "("
				duplicatesStr := make([]string, len(duplicates))
				for i, dup := range duplicates {
					duplicatesStr[i] = fmt.Sprintf("%s %s", t.yearStr(dup.year), dup.semester)
				}
				warning += strings.Join(duplicatesStr, ", ") + ")."
				winter1[ci1].warnings = append(winter1[ci1].warnings, warning)
			}
		}

		summer1 := year1.summer.courses
		for ci1 := range summer1 {
			duplicates := make([]struct {
				year     int
				semester string
			}, 0)
			for y2, year2 := range bp.years {
				winter2 := year2.winter.courses
				for ci2 := range winter2 {
					if summer1[ci1].code == winter2[ci2].code {
						duplicates = append(duplicates, struct {
							year     int
							semester string
						}{year: y2, semester: t.winter})
					}
				}
				summer2 := year2.summer.courses
				for ci2 := range summer2 {
					if summer1[ci1].code == summer2[ci2].code {
						duplicates = append(duplicates, struct {
							year     int
							semester string
						}{year: y2, semester: t.summer})
					}
				}
			}
			if len(duplicates) > 1 {
				warning := t.wAssignedMoreThanOnce + "("
				duplicatesStr := make([]string, len(duplicates))
				for i, dup := range duplicates {
					duplicatesStr[i] = fmt.Sprintf("%s %s", t.yearStr(dup.year), dup.semester)
				}
				warning += strings.Join(duplicatesStr, ", ") + ")."
				summer1[ci1].warnings = append(summer1[ci1].warnings, warning)
			}
		}
	}
	// TODO: add a warning if the course is not taught
}
