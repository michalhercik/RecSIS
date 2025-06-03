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

type Server struct {
	router *http.ServeMux
	Data   DataManager
	Auth   Authentication
	Page   Page
}

type Authentication interface {
	UserID(r *http.Request) (string, error)
}

type Page interface {
	View(main templ.Component, lang language.Language, title string, userID string) templ.Component
}

func (s *Server) Init() {
	s.router = http.NewServeMux()
	s.router.HandleFunc("GET /", s.page)
	s.router.HandleFunc("POST /year", s.yearAddition)
	s.router.HandleFunc("DELETE /year", s.yearRemoval)
	s.router.HandleFunc("PATCH /{year}/{semester}", s.foldSemester)
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

func generateWarnings(bp *Blueprint, t text) {
	// check if the course is assigned in a correct semester
	for _, year := range bp.years {
		winter := year.winter
		if !winter.folded {
			courses := winter.courses
			for ci := range courses {
				if courses[ci].Start == teachingSummerOnly {
					courses[ci].Warnings = append(courses[ci].Warnings, t.WWrongAssignWinter)
				}
			}
		}
		summer := year.summer
		if !summer.folded {
			courses := summer.courses
			for ci := range courses {
				if courses[ci].Start == teachingWinterOnly {
					courses[ci].Warnings = append(courses[ci].Warnings, t.WWrongAssignSummer)
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
					if winter1[ci1].Code == winter2[ci2].Code {
						duplicates = append(duplicates, struct {
							year     int
							semester string
						}{year: y2, semester: t.Winter})
					}
				}
				summer2 := year2.summer.courses
				for ci2 := range summer2 {
					if winter1[ci1].Code == summer2[ci2].Code {
						duplicates = append(duplicates, struct {
							year     int
							semester string
						}{year: y2, semester: t.Summer})
					}
				}
			}
			if len(duplicates) > 1 {
				warning := t.WAssignedMoreThanOnce + "("
				duplicatesStr := make([]string, len(duplicates))
				for i, dup := range duplicates {
					duplicatesStr[i] = fmt.Sprintf("%s %s", t.YearStr(dup.year), dup.semester)
				}
				warning += strings.Join(duplicatesStr, ", ") + ")."
				winter1[ci1].Warnings = append(winter1[ci1].Warnings, warning)
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
					if summer1[ci1].Code == winter2[ci2].Code {
						duplicates = append(duplicates, struct {
							year     int
							semester string
						}{year: y2, semester: t.Winter})
					}
				}
				summer2 := year2.summer.courses
				for ci2 := range summer2 {
					if summer1[ci1].Code == summer2[ci2].Code {
						duplicates = append(duplicates, struct {
							year     int
							semester string
						}{year: y2, semester: t.Summer})
					}
				}
			}
			if len(duplicates) > 1 {
				warning := t.WAssignedMoreThanOnce + "("
				duplicatesStr := make([]string, len(duplicates))
				for i, dup := range duplicates {
					duplicatesStr[i] = fmt.Sprintf("%s %s", t.YearStr(dup.year), dup.semester)
				}
				warning += strings.Join(duplicatesStr, ", ") + ")."
				summer1[ci1].Warnings = append(summer1[ci1].Warnings, warning)
			}
		}
	}
	// TODO: add a warning if the course is not taught
}

// ===============================================================================================================================
// Page
// ===============================================================================================================================

func (s Server) renderBlueprint(w http.ResponseWriter, r *http.Request, t text) {
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	data, err := s.Data.Blueprint(userID, t.Utils.Language)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	generateWarnings(data, t)
	Content(data, t).Render(r.Context(), w)
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
	data, err := s.Data.Blueprint(userID, t.Utils.Language)
	if err == nil {
		generateWarnings(data, t)
	}
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		result = Content(data, t)
	}
	s.Page.View(result, lang, t.PageTitle, userID).Render(r.Context(), w)
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

// type courseAdditionPresenter func(insertedCourseInfo) templ.Component

// func (s Server) courseAddition(w http.ResponseWriter, r *http.Request) {
// 	userID, err := s.Auth.UserID(r)
// 	if err != nil {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}
// 	course := r.PathValue("code")
// 	year, semester, err := parseYearSemester(r)
// 	if err != nil {
// 		log.Printf("courseAddition error: %v", err)
// 		return
// 	}
// 	reqSourceInt, err := strconv.Atoi(r.FormValue("requestSource"))
// 	if err != nil {
// 		reqSourceInt = int(sourceNone)
// 		log.Printf("courseAddition warning: %v", err)
// 	}
// 	reqSource := courseAdditionRequestSource(reqSourceInt)
// 	var presenter courseAdditionPresenter = DefaultCourseAdditionPresenter
// 	switch reqSource {
// 	case sourceBlueprint:
// 		// TODO: implement
// 		//
// 		// Example:
// 		//	file: blueprint/blueprint.templ
// 		//		templ Ribbon(insertInfo insertedCourseInfo) {
// 		//			<div class="ribbon"> {{insertInfo.courseID}} </div>
// 		//  	}
// 		//
// 		//	file: blueprint/blueprint.go
// 		//		presenter = Ribbon
// 	case sourceCourseDetail:
// 		// TODO: implement
// 	case sourceDegreePlan:
// 		// TODO: implement
// 	}
// 	// TODO: use this
// 	// text, err := parseLanguage(r)
// 	// if err != nil {
// 	// 	log.Println(err)
// 	// 	w.WriteHeader(http.StatusInternalServerError)
// 	// 	return
// 	// }
// 	courseID, err := s.Data.NewCourse(userID, course, year, semester)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	insertInfo := insertedCourseInfo{courseID: courseID, academicYear: year, semester: semester}
// 	presenter(insertInfo).Render(r.Context(), w)
// }

// ===============================================================================================================================
// Remove Year
// ===============================================================================================================================

func (s Server) yearRemoval(w http.ResponseWriter, r *http.Request) {
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
	unassign := r.FormValue("unassign")
	shouldUnassign, err := strconv.ParseBool(unassign)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if shouldUnassign {
		err = s.unassignYear(r, userID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
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
	err = s.Data.FoldSemester(userID, year, semester, folded)
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
