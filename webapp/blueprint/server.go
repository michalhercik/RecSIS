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

//================================================================================
// Server Type
//================================================================================

type Server struct {
	Auth   Authentication
	Data   Adapter
	Error  Error
	Page   Page
	router http.Handler
}

type Authentication interface {
	// Returns the user ID from an HTTP request.
	UserID(r *http.Request) string
}

type Error interface {
	// Logs the provided error.
	Log(err error)

	// Renders an error message to the user as a floating window, with a status code and localized message.
	Render(w http.ResponseWriter, r *http.Request, code int, userMsg string, lang language.Language)

	// Renders a full error page, including title and user ID, for major errors or page-level failures.
	RenderPage(w http.ResponseWriter, r *http.Request, code int, userMsg string, title string, userID string, lang language.Language)

	// Renders a fallback error page when a regular page cannot be rendered due to an error.
	CannotRenderPage(w http.ResponseWriter, r *http.Request, title string, userID string, err error, lang language.Language)
}

type Page interface {
	// Returns the page view component with injected main content, parameterized by language, title, and user ID.
	// Page adds header with navbar and footer.
	View(main templ.Component, lang language.Language, title string, userID string) templ.Component
}

//================================================================================
// Routing
//================================================================================

func (s Server) Router() http.Handler {
	return s.router
}

func (s *Server) Init() {
	router := http.NewServeMux()
	router.HandleFunc("GET /{$}", s.page)
	router.HandleFunc(fmt.Sprintf("PATCH /course/{%s}", recordID), s.courseMovement)
	router.HandleFunc("PATCH /courses", s.coursesMovement)
	router.HandleFunc(fmt.Sprintf("DELETE /course/{%s}", recordID), s.courseRemoval)
	router.HandleFunc("DELETE /courses", s.coursesRemoval)
	router.HandleFunc("POST /year", s.yearAddition)
	router.HandleFunc("DELETE /year", s.yearRemoval)
	router.HandleFunc("PATCH /fold", s.foldSemester)
	router.HandleFunc("/", s.pageNotFound)
	s.router = router
}

//================================================================================
// Page
//================================================================================

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID := s.Auth.UserID(r)
	data, err := s.Data.blueprint(userID, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, userID, lang)
		return
	}
	generateWarnings(data, t)
	main := Content(data, t)
	page := s.Page.View(main, lang, t.pageTitle, userID)
	err = page.Render(r.Context(), w)
	if err != nil {
		s.Error.CannotRenderPage(w, r, t.pageTitle, userID, errorx.AddContext(err), lang)
	}
}

func (s Server) pageNotFound(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID := s.Auth.UserID(r)
	s.Error.RenderPage(w, r, http.StatusNotFound, t.errPageNotFound, t.pageTitle, userID, lang)
}

//================================================================================
// Move Courses
//================================================================================

func (s Server) courseMovement(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID := s.Auth.UserID(r)
	courseID := r.PathValue(recordID)
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
	s.renderBlueprintContent(w, r, userID, lang)
}

func (s Server) coursesMovement(w http.ResponseWriter, r *http.Request) {
	var err error
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID := s.Auth.UserID(r)
	switch r.FormValue(typeParam) {
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
	userID := s.Auth.UserID(r)
	courseID := r.PathValue(recordID)
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
	s.renderBlueprintContent(w, r, userID, lang)
}

func (s Server) coursesRemoval(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID := s.Auth.UserID(r)
	var err error
	switch r.FormValue(typeParam) {
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
	userID := s.Auth.UserID(r)
	if err := s.Data.addYear(userID, lang); err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	s.renderBlueprintContent(w, r, userID, lang)
}

func (s Server) yearRemoval(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID := s.Auth.UserID(r)
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
	if err := s.Data.removeYear(userID, lang, shouldUnassign); err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	s.renderBlueprintContent(w, r, userID, lang)
}

//================================================================================
// Fold Semester
//================================================================================

func (s Server) foldSemester(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID := s.Auth.UserID(r)
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
	s.renderBlueprintContent(w, r, userID, lang)
}

//================================================================================
// Parse parameters
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

//================================================================================
// Generate Warnings
//================================================================================

func generateWarnings(bp *blueprintPage, t text) {
	forCorrectAssignment(bp, t)
	forDuplicateAssignments(bp, t)
	forRequisites(bp, t)
}

func forCorrectAssignment(bp *blueprintPage, t text) {
	for c := range bp.assignedCourses() {
		if c.course.semester == teachingWinterOnly && c.semester == assignmentSummer {
			c.course.warnings = append(c.course.warnings, t.wWrongAssignSummer)
		}
		if c.course.semester == teachingSummerOnly && c.semester == assignmentWinter {
			c.course.warnings = append(c.course.warnings, t.wWrongAssignWinter)
		}
	}
}

func forDuplicateAssignments(bp *blueprintPage, t text) {
	assignedCourseCodes := buildAssignedCourseMap(bp)

	for course := range bp.unassignedCourses() {
		if locations, exists := assignedCourseCodes[course.code]; exists {
			warning := makeDuplicateWarning(locations, t.wUnassignedButAssigned, t)
			course.warnings = append(course.warnings, warning)
		}
	}

	for _, locations := range assignedCourseCodes {
		if len(locations) > 1 {
			warning := makeDuplicateWarning(locations, t.wAssignedMoreThanOnce, t)
			for _, loc := range locations {
				loc.course.warnings = append(loc.course.warnings, warning)
			}
		}
	}
}

func makeDuplicateWarning(locations []courseLocation, warningText string, t text) string {
	warning := warningText + " ("
	locationsStr := make([]string, len(locations))
	for i, loc := range locations {
		semester := ""
		switch loc.semester {
		case assignmentWinter:
			semester = t.winter
		case assignmentSummer:
			semester = t.summer
		}
		locationsStr[i] = fmt.Sprintf("%s %s", t.yearStr(loc.year), semester)
	}
	warning += strings.Join(locationsStr, ", ") + ")."
	return warning
}

func forRequisites(bp *blueprintPage, t text) {
	assignedCourseMap := buildAssignedCourseMap(bp)

	// Check each course's requisites
	for courseLoc := range bp.assignedCourses() {
		course := courseLoc.course

		reqSlice := course.prerequisites
		if !reqSlice.isEmpty() {
			reqIsMet := areAllRequisitesMet(reqSlice, courseLoc, prerequisiteCondition, assignedCourseMap)
			if !reqIsMet {
				course.warnings = append(course.warnings, t.wPrerequisiteNotMet)
				// TODO: better warning message with details
			}
		}

		reqSlice = course.corequisites
		if !reqSlice.isEmpty() {
			reqIsMet := areAllRequisitesMet(reqSlice, courseLoc, corequisiteCondition, assignedCourseMap)
			if !reqIsMet {
				course.warnings = append(course.warnings, t.wCorequisiteNotMet)
				// TODO: better warning message with details
			}
		}

		reqSlice = course.incompatibles
		if !reqSlice.isEmpty() {
			reqIsViolated := isAnyRequisiteMet(reqSlice, courseLoc, incompatibleCondition, assignedCourseMap)
			if reqIsViolated {
				course.warnings = append(course.warnings, t.wIncompatiblePresent)
				// TODO: better warning message with details
			}
		}
	}
}

func areAllRequisitesMet(courses []requisite, originLoc courseLocation, condition requisiteCondition, courseMap map[string][]courseLocation) bool {
	for _, req := range courses {
		if req.isNode() {
			if !isCourseRequisiteMet(req, originLoc, condition, courseMap) {
				return false
			}
		} else if req.isDisjunction() {
			// At least one child requisite must be met
			if !isAnyRequisiteMet(req.children, originLoc, condition, courseMap) {
				return false
			}
		} else if req.isConjunction() {
			// All child requisites must be met
			if !areAllRequisitesMet(req.children, originLoc, condition, courseMap) {
				return false
			}
		}
	}
	return true
}

func isAnyRequisiteMet(courses []requisite, originLoc courseLocation, condition requisiteCondition, courseMap map[string][]courseLocation) bool {
	for _, req := range courses {
		if req.isNode() {
			if isCourseRequisiteMet(req, originLoc, condition, courseMap) {
				return true
			}
		} else if req.isDisjunction() {
			// At least one child requisite must be met
			if isAnyRequisiteMet(req.children, originLoc, condition, courseMap) {
				return true
			}
		} else if req.isConjunction() {
			// All child requisites must be met
			if areAllRequisitesMet(req.children, originLoc, condition, courseMap) {
				return true
			}
		}
	}
	return false
}

func isCourseRequisiteMet(reqNode requisite, originLoc courseLocation, condition requisiteCondition, courseMap map[string][]courseLocation) bool {
	if locs, exists := courseMap[reqNode.courseCode]; exists {
		for _, reqLoc := range locs {
			if condition(originLoc, reqLoc) {
				return true
			}
		}
	}
	return false
}

// Helper: Build map of assigned courses' codes to their locations
func buildAssignedCourseMap(bp *blueprintPage) map[string][]courseLocation {
	courseMap := make(map[string][]courseLocation)
	for loc := range bp.assignedCourses() {
		courseMap[loc.course.code] = append(courseMap[loc.course.code], loc)
	}
	return courseMap
}

//================================================================================
// Render Content
//================================================================================

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
	view := Content(data, t)
	err = view.Render(r.Context(), w)
	if err != nil {
		s.Error.CannotRenderPage(w, r, t.pageTitle, userID, errorx.AddContext(err), lang)
	}
}
