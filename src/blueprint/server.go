package blueprint

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/language"
)

/**
TODO:
	- document functions
*/

//================================================================================
// Server Type
//================================================================================

type Server struct {
	Auth   Authentication
	Data   DBManager
	Error  Error
	Page   Page
	router http.Handler
}

func (s Server) Router() http.Handler {
	return s.router
}

type Authentication interface {
	UserID(r *http.Request) (string, error)
}

type Error interface {
	Log(err error)
	Render(w http.ResponseWriter, r *http.Request, code int, userMsg string, lang language.Language)
	RenderPage(w http.ResponseWriter, r *http.Request, code int, userMsg string, title string, userID string, lang language.Language)
	CannotRenderPage(w http.ResponseWriter, r *http.Request, title string, userID string, err error, lang language.Language)
}

type Page interface {
	View(main templ.Component, lang language.Language, title string, userID string) templ.Component
}

//================================================================================
// Routing
//================================================================================

func (s *Server) Init() {
	router := http.NewServeMux()
	router.HandleFunc("GET /{$}", s.page)
	router.HandleFunc("PATCH /course/{id}", s.courseMovement)
	router.HandleFunc("PATCH /courses", s.coursesMovement)
	router.HandleFunc("DELETE /course/{id}", s.courseRemoval)
	router.HandleFunc("DELETE /courses", s.coursesRemoval)
	router.HandleFunc("POST /year", s.yearAddition)
	router.HandleFunc("DELETE /year", s.yearRemoval)
	router.HandleFunc("PATCH /fold", s.foldSemester)

	// Wrap mux to catch unmatched routes
	s.router = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if mux has a handler for the URL
		_, pattern := router.Handler(r)
		if pattern == "" {
			s.pageNotFound(w, r)
			return
		}
		router.ServeHTTP(w, r)
	})
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
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, "", lang)
		return
	}
	data, err := s.Data.blueprint(userID, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, userID, lang)
		return
	}
	generateWarnings(data, t)
	result = Content(data, t)
	err = s.Page.View(result, lang, t.pageTitle, userID).Render(r.Context(), w)
	if err != nil {
		s.Error.CannotRenderPage(w, r, t.pageTitle, userID, errorx.AddContext(err), lang)
	}
}

func (s Server) renderBlueprintContent(w http.ResponseWriter, r *http.Request, userID string, lang language.Language) {
	t := texts[lang]
	data, err := s.Data.blueprint(userID, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	generateWarnings(data, t)
	err = Content(data, t).Render(r.Context(), w)
	if err != nil {
		s.Error.CannotRenderPage(w, r, t.pageTitle, userID, errorx.AddContext(err), lang)
	}
}

//================================================================================
// Move Courses
//================================================================================

func (s Server) courseMovement(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID, err := s.Auth.UserID(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, "", lang)
		return
	}
	courseID := r.PathValue("id")
	if courseID == "" {
		s.Error.Log(errorx.AddContext(fmt.Errorf("course ID is missing in the request path")))
		s.Error.Render(w, r, http.StatusBadRequest, t.errMissingCourseID, lang)
		return
	}
	course, err := strconv.Atoi(courseID)
	if err != nil {
		s.Error.Log(errorx.AddContext(fmt.Errorf("unable to parse course ID to int: %w", err)))
		s.Error.Render(w, r, http.StatusBadRequest, t.errInvalidCourseID, lang)
	}
	year, semester, position, err := parseYearSemesterPosition(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
	}
	if position == lastPosition {
		err = s.Data.appendCourses(userID, lang, year, semester, course)
	} else if position >= 0 {
		err = s.Data.moveCourses(userID, lang, year, semester, position, course)
	} else {
		s.Error.Log(errorx.AddContext(fmt.Errorf("invalid position %d", position)))
		s.Error.Render(w, r, http.StatusBadRequest, t.errInvalidPositionParam, lang)
		return
	}
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	// TODO: Render just the necessary for performance
	s.renderBlueprintContent(w, r, userID, lang)
}

func (s Server) coursesMovement(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID, err := s.Auth.UserID(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, "", lang)
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
		s.Error.Log(errorx.AddContext(fmt.Errorf("invalid type %s", r.FormValue(typeParam))))
		s.Error.Render(w, r, http.StatusBadRequest, t.errInvalidMoveType, lang)
		return
	}
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	s.renderBlueprintContent(w, r, userID, lang)
}

func (s Server) unassignYear(r *http.Request, userID string) error {
	lang := language.FromContext(r.Context())
	year, err := parseYear(r)
	if err != nil {
		return errorx.AddContext(err)
	}
	err = s.Data.unassignYear(userID, lang, year)
	if err != nil {
		return errorx.AddContext(err)
	}
	return nil
}

func (s Server) unassignSemester(r *http.Request, userID string) error {
	lang := language.FromContext(r.Context())
	year, semester, err := parseYearSemester(r)
	if err != nil {
		return errorx.AddContext(err)
	}
	err = s.Data.unassignSemester(userID, lang, year, semester)
	if err != nil {
		return errorx.AddContext(err)
	}
	return nil
}

func (s Server) moveCourses(r *http.Request, userID string) error {
	lang := language.FromContext(r.Context())
	year, semester, position, err := parseYearSemesterPosition(r)
	if err != nil {
		return errorx.AddContext(err)
	}
	courses, err := atoiSliceCourses(r.Form[checkboxName], lang)
	if err != nil {
		return errorx.AddContext(err)
	}
	if position == lastPosition {
		err = s.Data.appendCourses(userID, lang, year, semester, courses...)
	} else if position > 0 {
		err = s.Data.moveCourses(userID, lang, year, semester, position, courses...)
	} else {
		return errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("invalid position %d", position)),
			http.StatusBadRequest,
			texts[lang].errInvalidPositionParam,
		)
	}
	return err
}

// ================================================================================
// Remove Courses
// ================================================================================

func (s Server) courseRemoval(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID, err := s.Auth.UserID(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, "", lang)
		return
	}
	courseID := r.PathValue("id")
	if courseID == "" {
		s.Error.Log(errorx.AddContext(fmt.Errorf("course ID is missing in the request path")))
		s.Error.Render(w, r, http.StatusBadRequest, t.errMissingCourseID, lang)
		return
	}
	course, err := strconv.Atoi(courseID)
	if err != nil {
		s.Error.Log(errorx.AddContext(fmt.Errorf("unable to parse course ID to int: %w", err)))
		s.Error.Render(w, r, http.StatusBadRequest, t.errInvalidCourseID, lang)
		return
	}
	err = s.Data.removeCourses(userID, lang, course)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	// TODO: Render just credits stats with own sql query for performance
	s.renderBlueprintContent(w, r, userID, lang)
}

func (s Server) coursesRemoval(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID, err := s.Auth.UserID(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, "", lang)
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
		s.Error.Log(errorx.AddContext(fmt.Errorf("invalid type %s", r.FormValue(typeParam))))
		s.Error.Render(w, r, http.StatusBadRequest, t.errInvalidRemoveType, lang)
		return
	}
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	s.renderBlueprintContent(w, r, userID, lang)
}

func (s Server) removeCoursesByYear(r *http.Request, userID string) error {
	lang := language.FromContext(r.Context())
	year, err := parseYear(r)
	if err != nil {
		return errorx.AddContext(err)
	}
	err = s.Data.removeCoursesByYear(userID, lang, year)
	if err != nil {
		return errorx.AddContext(err)
	}
	return nil
}

func (s Server) removeCoursesBySemester(r *http.Request, userID string) error {
	lang := language.FromContext(r.Context())
	year, semester, err := parseYearSemester(r)
	if err != nil {
		return errorx.AddContext(err)
	}
	err = s.Data.removeCoursesBySemester(userID, lang, year, semester)
	if err != nil {
		return errorx.AddContext(err)
	}
	return nil
}

func (s Server) removeCourses(r *http.Request, userID string) error {
	lang := language.FromContext(r.Context())
	r.ParseForm()
	courses, err := atoiSliceCourses(r.Form[checkboxName], lang)
	if err != nil {
		return errorx.AddContext(err)
	}
	err = s.Data.removeCourses(userID, lang, courses...)
	if err != nil {
		return errorx.AddContext(err)
	}
	return nil
}

//================================================================================
// Add/Remove Year
//================================================================================

func (s Server) yearAddition(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID, err := s.Auth.UserID(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, "", lang)
		return
	}
	if err := s.Data.addYear(userID, lang); err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	// TODO: Render just credits stats with own sql query for performance
	s.renderBlueprintContent(w, r, userID, lang)
}

func (s Server) yearRemoval(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID, err := s.Auth.UserID(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, "", lang)
		return
	}
	unassign := r.FormValue(unassignParam)
	if unassign == "" {
		s.Error.Log(errorx.AddContext(fmt.Errorf("unassign parameter is missing in the request form")))
		s.Error.Render(w, r, http.StatusBadRequest, t.errMissingUnassignParam, lang)
		return
	}
	shouldUnassign, err := strconv.ParseBool(unassign)
	if err != nil {
		s.Error.Log(errorx.AddContext(fmt.Errorf("unable to parse unassign parameter to bool: %w", err)))
		s.Error.Render(w, r, http.StatusBadRequest, t.errInvalidUnassignParam, lang)
		return
	}
	year, err := parseYear(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	if err := s.Data.removeYear(userID, lang, year, shouldUnassign); err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	// TODO: Render just credits stats with own sql query for performance
	s.renderBlueprintContent(w, r, userID, lang)
}

//================================================================================
// Fold Semester
//================================================================================

func (s Server) foldSemester(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID, err := s.Auth.UserID(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, "", lang)
		return
	}
	year, semester, err := parseYearSemester(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	foldedString := r.FormValue(foldedParam)
	if foldedString == "" {
		s.Error.Log(errorx.AddContext(fmt.Errorf("folded parameter is missing in the request form")))
		s.Error.Render(w, r, http.StatusBadRequest, t.errMissingFoldedParam, lang)
		return
	}
	folded, err := strconv.ParseBool(foldedString)
	if err != nil {
		s.Error.Log(errorx.AddContext(fmt.Errorf("unable to parse folded parameter to bool: %w", err)))
		s.Error.Render(w, r, http.StatusBadRequest, t.errInvalidFoldedParam, lang)
		return
	}
	err = s.Data.foldSemester(userID, lang, year, semester, folded)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	// TODO: Render just targeted table for performance
	s.renderBlueprintContent(w, r, userID, lang)
}

//================================================================================
// Utils
//================================================================================

func parseYear(r *http.Request) (int, error) {
	lang := language.FromContext(r.Context())
	t := texts[lang]

	yearString := r.FormValue(yearParam)
	if yearString == "" {
		return 0, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("year parameter is missing in the request form")),
			http.StatusBadRequest,
			t.errMissingYearParam,
		)
	}

	year, err := strconv.Atoi(yearString)
	if err != nil {
		return 0, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("unable to parse year parameter to int: %w", err)),
			http.StatusBadRequest,
			t.errInvalidYearParam,
		)
	}

	if year < 0 {
		return 0, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("invalid year %d", year)),
			http.StatusBadRequest,
			t.errInvalidYearParam,
		)
	}

	return year, nil
}

func parseSemester(r *http.Request) (semesterAssignment, error) {
	lang := language.FromContext(r.Context())
	t := texts[lang]

	semesterString := r.FormValue(semesterParam)
	if semesterString == "" {
		return semesterAssignment(0), errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("semester parameter is missing in the request form")),
			http.StatusBadRequest,
			t.errMissingSemesterParam,
		)
	}

	semesterInt, err := strconv.Atoi(semesterString)
	if err != nil {
		return semesterAssignment(0), errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("unable to parse semester parameter to int: %w", err)),
			http.StatusBadRequest,
			t.errInvalidSemesterParam,
		)
	}

	if semesterInt < int(assignmentNone) || semesterInt > int(assignmentSummer) {
		return semesterAssignment(0), errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("invalid semester %d", semesterInt)),
			http.StatusBadRequest,
			t.errInvalidSemesterParam,
		)
	}

	semester := semesterAssignment(semesterInt)
	return semester, nil
}

func parseYearSemester(r *http.Request) (int, semesterAssignment, error) {
	year, err := parseYear(r)
	if err != nil {
		return 0, semesterAssignment(0), errorx.AddContext(err)
	}
	semester, err := parseSemester(r)
	if err != nil {
		return year, semesterAssignment(0), errorx.AddContext(err)
	}
	return year, semester, nil
}

func parsePosition(r *http.Request) (int, error) {
	lang := language.FromContext(r.Context())
	t := texts[lang]

	positionString := r.FormValue(positionParam)
	if positionString == "" {
		return 0, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("position parameter is missing in the request form")),
			http.StatusBadRequest,
			t.errMissingPositionParam,
		)
	}

	position, err := strconv.Atoi(positionString)
	if err != nil {
		return 0, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("unable to parse position parameter to int: %w", err)),
			http.StatusBadRequest,
			t.errInvalidPositionParam,
		)
	}

	return position, nil
}

func parseYearSemesterPosition(r *http.Request) (int, semesterAssignment, int, error) {
	year, semester, err := parseYearSemester(r)
	if err != nil {
		return year, semester, 0, errorx.AddContext(err)
	}

	position, err := parsePosition(r)
	if err != nil {
		return year, semester, 0, errorx.AddContext(err)
	}

	return year, semester, position, nil
}

func atoiSliceCourses(s []string, lang language.Language) ([]int, error) {
	result := make([]int, len(s))
	for i, elem := range s {
		var err error
		result[i], err = strconv.Atoi(elem)
		if err != nil {
			return nil, errorx.NewHTTPErr(
				errorx.AddContext(fmt.Errorf("unable to parse string %s to int: %w", elem, err)),
				http.StatusBadRequest,
				texts[lang].errInvalidCourseID,
			)
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

//================================================================================
// Page Not Found
//================================================================================

func (s Server) pageNotFound(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID, err := s.Auth.UserID(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, "", lang)
		return
	}
	s.Error.RenderPage(w, r, http.StatusNotFound, t.errPageNotFound, t.pageTitle, userID, lang)
}
