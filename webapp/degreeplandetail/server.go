package degreeplandetail

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
	Auth                Authentication
	BpBtn               BlueprintAddButton
	Data                DBManager
	Error               Error
	SearchRedirectPath  string
	ComparePlanUrlParam string
	Page                Page
	router              http.Handler
}

type Authentication interface {
	// Returns the user ID from an HTTP request.
	UserID(r *http.Request) string
}

type BlueprintAddButton interface {
	// Returns a partial component for rendering add button on every row, parameterized by language.
	PartialComponent(lang language.Language) PartialBlueprintAdd

	// Returns a partial component for rendering add button for checked courses, parameterized by language.
	PartialComponentSecond(lang language.Language) PartialBlueprintAdd

	// Parses the HTTP request to extract course codes and blueprint context (year, semester).
	ParseRequest(r *http.Request, additionalCourses []string) ([]string, int, int, error)

	// Executes the action of adding one or more courses to the blueprint for a user.
	Action(userID string, year int, semester int, lang language.Language, course ...string) ([]int, error)

	// Returns the endpoint (method and URL path string) for the add-to-blueprint action.
	Endpoint() string
}

type PartialBlueprintAdd = func(hxSwap, hxTarget, hxInclude string, years []bool, course string) templ.Component

type Error interface {
	// Logs the provided error.
	Log(err error)

	// Renders an error message to the user as a floating window, with a status code and localized message.
	Render(w http.ResponseWriter, r *http.Request, code int, userMsg string, lang language.Language)

	// Renders a full error page, including title and user ID, for major errors or page-level failures.
	RenderPage(w http.ResponseWriter, r *http.Request, code int, userMsg string, title string, userID string, lang language.Language)

	// Renders a fallback error page when a regular page cannot be rendered due to an error.
	CannotRenderPage(w http.ResponseWriter, r *http.Request, title string, userID string, err error, lang language.Language)

	// Renders a floating window with error when any component cannot be rendered due to an error.
	CannotRenderComponent(w http.ResponseWriter, r *http.Request, err error, lang language.Language)
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
	type routingElement struct {
		path    string
		handler http.HandlerFunc
		params  []any
	}
	endpoints := []routingElement{
		{"GET /{$}", s.userDegreePlanPage, nil},
		{"GET /{%s}", s.degreePlanByCodePage, []any{dpCode}},
		{"PATCH /{%s}", s.saveDegreePlan, []any{dpCode}},
		{"DELETE /", s.deleteSavedPlan, nil},
		{"%s", s.addCourseToBlueprint, []any{s.BpBtn.Endpoint()}},
		{"PATCH /plan-to-blueprint/{%s}", s.mergeRecPlanWithBlueprint, []any{dpCode}},
		{"PUT /plan-to-blueprint/{%s}", s.rewriteBlueprintWithRecPlan, []any{dpCode}},
		{"/", s.pageNotFound, nil},
	}
	router := http.NewServeMux()
	for _, e := range endpoints {
		router.HandleFunc(fmt.Sprintf(e.path, e.params...), e.handler)
	}
	s.router = router
}

//================================================================================
// Handlers
//================================================================================

func (s Server) userDegreePlanPage(w http.ResponseWriter, r *http.Request) {
	userID := s.Auth.UserID(r)
	lang := language.FromContext(r.Context())
	t := texts[lang]
	if !s.Data.userHasSelectedDegreePlan(userID) {
		http.Redirect(w, r, lang.LocalizeURL(s.SearchRedirectPath), http.StatusSeeOther)
		return
	}
	dp, err := s.Data.userDegreePlan(userID, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, userID, lang)
		return
	}
	main := s.pageContent(dp, t)
	page := s.Page.View(main, lang, t.pageTitle, userID)
	err = page.Render(r.Context(), w)
	if err != nil {
		s.Error.CannotRenderPage(w, r, t.pageTitle, userID, errorx.AddContext(err), lang)
	}
}

func (s Server) degreePlanByCodePage(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID := s.Auth.UserID(r)
	dpCode := r.PathValue(dpCode)
	dp, err := s.Data.degreePlan(userID, dpCode, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.pageTitle, userID, lang)
		return
	}
	main := s.pageContent(dp, t)
	page := s.Page.View(main, lang, t.pageTitle, userID)
	err = page.Render(r.Context(), w)
	if err != nil {
		s.Error.CannotRenderPage(w, r, t.pageTitle, userID, errorx.AddContext(err), lang)
	}
}

func (s Server) saveDegreePlan(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	userID := s.Auth.UserID(r)
	dpCode := r.PathValue(dpCode)
	err := s.Data.saveDegreePlan(userID, dpCode, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	s.userDegreePlanPage(w, r)
}

func (s Server) deleteSavedPlan(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	userID := s.Auth.UserID(r)
	planCode, err := s.Data.deleteSavedDegreePlan(userID, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	r.SetPathValue(dpCode, planCode)
	s.degreePlanByCodePage(w, r)
}

func (s Server) addCourseToBlueprint(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	userID := s.Auth.UserID(r)
	r.ParseForm()
	// to include courses selected by checkboxes
	additionalCourses := r.Form[checkboxName]
	courseCodes, year, semester, err := s.BpBtn.ParseRequest(r, additionalCourses)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	_, err = s.BpBtn.Action(userID, year, semester, lang, courseCodes...)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	dp, err := s.getCorrectPlanPage(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	t := texts[lang]
	view := s.pageContent(dp, t)
	err = view.Render(r.Context(), w)
	if err != nil {
		s.Error.CannotRenderComponent(w, r, errorx.AddContext(err), lang)
	}
}

func (s Server) mergeRecPlanWithBlueprint(w http.ResponseWriter, r *http.Request) {
	planCode := r.PathValue(dpCode)
	lang := language.FromContext(r.Context())
	userID := s.Auth.UserID(r)
	maxYear, err := parseMaxYearParam(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	err = s.Data.mergeRecommendedPlanWithBlueprint(userID, planCode, maxYear, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	dp, err := s.getCorrectPlanPage(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	t := texts[lang]
	view := s.pageContent(dp, t)
	err = view.Render(r.Context(), w)
	if err != nil {
		s.Error.CannotRenderComponent(w, r, errorx.AddContext(err), lang)
	}
}

func (s Server) rewriteBlueprintWithRecPlan(w http.ResponseWriter, r *http.Request) {
	planCode := r.PathValue(dpCode)
	lang := language.FromContext(r.Context())
	userID := s.Auth.UserID(r)
	maxYear, err := parseMaxYearParam(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	err = s.Data.rewriteBlueprintWithRecommendedPlan(userID, planCode, maxYear, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	dp, err := s.getCorrectPlanPage(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	t := texts[lang]
	view := s.pageContent(dp, t)
	err = view.Render(r.Context(), w)
	if err != nil {
		s.Error.CannotRenderComponent(w, r, errorx.AddContext(err), lang)
	}
}

/*
TODO: better handler this problem
  - do not rely on referer header
  - if there is no referer, redirect to user's saved plan may be wrong, as there may be no saved plan
  - perhaps pass the plan code as a query parameter when needed
  - although, this is only used after adding a course to blueprint from this page
    -- so referer should always be present
    -- and it should contain /degreeplan/{code} or /degreeplan/
*/
func (s Server) getCorrectPlanPage(r *http.Request) (*degreePlanPage, error) {
	lang := language.FromContext(r.Context())
	userID := s.Auth.UserID(r)
	var dp *degreePlanPage
	var err error
	parts := strings.Split(r.Referer(), "/degreeplan/")
	if len(parts) > 1 && parts[1] != "" {
		planCode := parts[1]
		dp, err = s.Data.degreePlan(userID, planCode, lang)
	} else {
		dp, err = s.Data.userDegreePlan(userID, lang)
	}
	return dp, err
}

func (s Server) pageContent(dp *degreePlanPage, t text) templ.Component {
	partialBpBtn := s.BpBtn.PartialComponent(t.language)
	partialBpBtnChecked := s.BpBtn.PartialComponentSecond(t.language)
	dp.searchEndpoint = s.SearchRedirectPath
	dp.compareParam = s.ComparePlanUrlParam
	main := Content(dp, t, partialBpBtn, partialBpBtnChecked)
	return main
}

func (s Server) pageNotFound(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID := s.Auth.UserID(r)
	s.Error.RenderPage(w, r, http.StatusNotFound, t.errPageNotFound, t.pageTitle, userID, lang)
}

//================================================================================
// Parse parameters
//================================================================================

func parseMaxYearParam(r *http.Request) (int, error) {
	lang := language.FromContext(r.Context())
	t := texts[lang]

	yearString := r.FormValue(maxYearParam)
	if yearString == "" {
		return 0, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("year parameter is missing in the request form")),
			http.StatusBadRequest,
			t.errMissingMaxYearParam,
		)
	}

	year, err := strconv.Atoi(yearString)
	if err != nil {
		return 0, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("unable to parse year parameter to int: %w", err)),
			http.StatusBadRequest,
			t.errInvalidMaxYearParam,
		)
	}

	if year < 0 {
		return 0, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("invalid year %d", year)),
			http.StatusBadRequest,
			t.errInvalidMaxYearParam,
		)
	}

	return year, nil
}
