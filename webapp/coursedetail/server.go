package coursedetail

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/filters"
	"github.com/michalhercik/RecSIS/language"
)

//================================================================================
// Server Type
//================================================================================

type Server struct {
	Auth    Authentication
	BpBtn   BlueprintAddButton
	Data    DBManager
	Error   Error
	Filters filters.Filters
	Search  Search
	Page    Page
	router  http.Handler
}

func (s *Server) Init() {
	var err error
	if err = s.Filters.Init(); err != nil {
		log.Fatal("coursedetail.Init: ", err)
	}
	s.initRouter()
}

type Authentication interface {
	UserID(r *http.Request) string
}

type BlueprintAddButton interface {
	PartialComponent(lang language.Language) PartialBlueprintAdd
	ParseRequest(r *http.Request, additionalCourses []string) ([]string, int, int, error)
	Action(userID string, year int, semester int, lang language.Language, course ...string) ([]int, error)
	Endpoint() string
}

type PartialBlueprintAdd = func(hxSwap, hxTarget, hxInclude string, semesters []bool, course string) templ.Component

type Error interface {
	Log(err error)
	Render(w http.ResponseWriter, r *http.Request, code int, userMsg string, lang language.Language)
	RenderPage(w http.ResponseWriter, r *http.Request, code int, userMsg string, title string, userID string, lang language.Language)
	CannotRenderPage(w http.ResponseWriter, r *http.Request, title string, userID string, err error, lang language.Language)
	CannotRenderComponent(w http.ResponseWriter, r *http.Request, err error, lang language.Language)
}

type Page interface {
	View(main templ.Component, lang language.Language, title string, userID string) templ.Component
}

//================================================================================
// Routing
//================================================================================

func (s Server) Router() http.Handler {
	return s.router
}

func (s *Server) initRouter() {
	type routingElement struct {
		path    string
		handler http.HandlerFunc
		params  []any
	}
	endpoints := []routingElement{
		{"GET /{%s}", s.page, []any{courseCode}},
		{"GET /survey/{%s}", s.survey, []any{courseCode}},
		{"GET /survey/next/{%s}", s.surveyNext, []any{courseCode}},
		{"PUT /rating/{%s}/{%s}", s.rateCategory, []any{courseCode, ratingCategory}},
		{"DELETE /rating/{%s}/{%s}", s.deleteCategoryRating, []any{courseCode, ratingCategory}},
		{"PUT /rating/{%s}", s.rate, []any{courseCode}},
		{"DELETE /rating/{%s}", s.deleteRating, []any{courseCode}},
		{"%s", s.addCourseToBlueprint, []any{s.BpBtn.Endpoint()}},
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

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID := s.Auth.UserID(r)
	code := r.PathValue(courseCode)
	course, err := s.Data.course(userID, code, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.errPageTitle, userID, lang)
		return
	}
	title := course.code + " - " + course.title
	btn := s.BpBtn.PartialComponent(lang)
	courseDetailPage := courseDetailPage{
		course: course,
	}
	main := Content(&courseDetailPage, t, btn)
	page := s.Page.View(main, lang, title, userID)
	err = page.Render(r.Context(), w)
	if err != nil {
		s.Error.CannotRenderPage(w, r, title, userID, errorx.AddContext(err), lang)
	}
}

func (s Server) survey(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	model, err := s.surveyViewModel(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, model.lang)
		userMsg = fmt.Sprintf("%s: %s", texts[lang].errCannotLoadSurvey, userMsg)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	content := SurveyFiltersContent(model, texts[model.lang])
	err = content.Render(r.Context(), w)
	if err != nil {
		s.Error.CannotRenderComponent(w, r, errorx.AddContext(err), model.lang)
	}
}

func (s Server) surveyNext(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	model, err := s.surveyViewModel(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, model.lang)
		userMsg = fmt.Sprintf("%s: %s", texts[lang].errCannotLoadSurvey, userMsg)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	content := SurveysContent(model, texts[model.lang])
	err = content.Render(r.Context(), w)
	if err != nil {
		s.Error.CannotRenderComponent(w, r, errorx.AddContext(err), model.lang)
	}
}

func (s Server) surveyViewModel(r *http.Request) (surveyViewModel, error) {
	req, err := s.parseQueryRequest(r)
	if err != nil {
		return surveyViewModel{}, errorx.AddContext(err)
	}
	code := r.PathValue(courseCode)
	req.filter.Append(meiliCourseCode, code)
	searchResponse, err := s.Search.comments(req)
	if err != nil {
		return surveyViewModel{}, errorx.AddContext(err, errorx.P(courseCode, code))
	}
	result := surveyViewModel{
		lang:   req.lang,
		code:   code,
		survey: searchResponse.Survey,
		offset: req.offset,
		isEnd:  searchResponse.EstimatedTotalHits <= req.offset+req.limit,
		facets: s.Filters.IterFiltersWithFacets(searchResponse.FacetDistribution, r.URL.Query(), req.lang),
		query:  req.query,
	}
	return result, nil
}

func (s Server) parseQueryRequest(r *http.Request) (request, error) {
	var req request
	userID := s.Auth.UserID(r)
	lang := language.FromContext(r.Context())
	query := r.FormValue(searchQuery)
	offset, err := strconv.Atoi(r.FormValue(surveyOffset))
	if err != nil {
		offset = defaultSurveyOffset
	}
	filter, err := s.Filters.ParseURLQuery(r.URL.Query(), lang)
	if err != nil {
		return req, errorx.AddContext(err)
	}

	req = request{
		userID: userID,
		query:  query,
		offset: offset,
		limit:  resultsPerPage,
		lang:   lang,
		filter: expression(&filter),
		facets: s.Filters.Facets(),
		sort:   meiliSort,
	}
	return req, nil
}

func (s Server) rateCategory(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	userID := s.Auth.UserID(r)
	code := r.PathValue(courseCode)
	category := r.PathValue(ratingCategory)
	ratingString := r.FormValue(ratingParam)
	rating, err := strconv.Atoi(ratingString)
	if err != nil {
		s.Error.Log(errorx.AddContext(err, errorx.P(ratingParam, ratingString)))
		s.Error.Render(w, r, http.StatusBadRequest, texts[lang].errRatingMustBeInt, lang)
		return
	}
	if rating < minRating || rating > maxRating {
		s.Error.Log(errorx.AddContext(fmt.Errorf("rating is not between %d and %d", minRating, maxRating), errorx.P(ratingParam, ratingString)))
		s.Error.Render(w, r, http.StatusBadRequest, texts[lang].errInvalidRating0to10, lang)
		return
	}
	updatedRating, err := s.Data.rateCategory(userID, code, category, rating, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	content := CategoryRating(updatedRating, code, texts[lang])
	err = content.Render(r.Context(), w)
	if err != nil {
		s.Error.CannotRenderComponent(w, r, errorx.AddContext(err), lang)
	}
}

func (s Server) deleteCategoryRating(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	userID := s.Auth.UserID(r)
	code := r.PathValue(courseCode)
	category := r.PathValue(ratingCategory)
	updatedRating, err := s.Data.deleteCategoryRating(userID, code, category, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	content := CategoryRating(updatedRating, code, texts[lang])
	err = content.Render(r.Context(), w)
	if err != nil {
		s.Error.CannotRenderComponent(w, r, errorx.AddContext(err), lang)
	}
}

func (s Server) rate(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	userID := s.Auth.UserID(r)
	code := r.PathValue(courseCode)
	ratingString := r.FormValue(ratingParam)
	rating, err := strconv.Atoi(ratingString)
	if err != nil {
		s.Error.Log(errorx.AddContext(err, errorx.P(ratingParam, ratingString)))
		s.Error.Render(w, r, http.StatusBadRequest, texts[lang].errRatingMustBeInt, lang)
		return
	}
	if rating != negativeRating && rating != positiveRating {
		s.Error.Log(errorx.AddContext(fmt.Errorf("rating is not %d or %d", negativeRating, positiveRating), errorx.P(ratingParam, ratingString)))
		s.Error.Render(w, r, http.StatusBadRequest, texts[lang].errInvalidRating0or1, lang)
		return
	}
	updatedRating, err := s.Data.rate(userID, code, rating, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	content := OverallRating(updatedRating, code, texts[lang])
	err = content.Render(r.Context(), w)
	if err != nil {
		s.Error.CannotRenderComponent(w, r, errorx.AddContext(err), lang)
	}
}

func (s Server) deleteRating(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	userID := s.Auth.UserID(r)
	code := r.PathValue(courseCode)
	updatedRating, err := s.Data.deleteRating(userID, code, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	content := OverallRating(updatedRating, code, texts[lang])
	err = content.Render(r.Context(), w)
	if err != nil {
		s.Error.CannotRenderComponent(w, r, errorx.AddContext(err), lang)
	}
}

func (s Server) addCourseToBlueprint(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID := s.Auth.UserID(r)
	courseCodes, year, semester, err := s.BpBtn.ParseRequest(r, nil)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	if len(courseCodes) != 1 {
		s.Error.Log(errorx.AddContext(fmt.Errorf("expected exactly one course code, got %d", len(courseCodes)), errorx.P("courseCodes", courseCodes)))
		s.Error.Render(w, r, http.StatusBadRequest, t.errUnexpectedNumberOfCourses, lang)
		return
	}
	courseCode := courseCodes[0]
	_, err = s.BpBtn.Action(userID, year, semester, lang, courseCode)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	btn := s.BpBtn.PartialComponent(lang)
	course, err := s.Data.course(userID, courseCode, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.errPageTitle, userID, lang)
		return
	}
	courseDetailPage := courseDetailPage{
		course: course,
	}
	content := Content(&courseDetailPage, t, btn)
	err = content.Render(r.Context(), w)
	if err != nil {
		s.Error.CannotRenderPage(w, r, t.errPageTitle, userID, errorx.AddContext(err), lang)
	}
}

func (s Server) pageNotFound(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID := s.Auth.UserID(r)
	s.Error.RenderPage(w, r, http.StatusNotFound, t.errPageNotFound, t.errPageTitle, userID, lang)
}
