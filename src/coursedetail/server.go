package coursedetail

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
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

func (s Server) Router() http.Handler {
	return s.router
}

type Authentication interface {
	UserID(r *http.Request) (string, error)
}

type BlueprintAddButton interface {
	PartialComponent(lang language.Language) PartialBlueprintAdd
	ParseRequest(r *http.Request, additionalCourses []string) ([]string, int, int, error)
	Action(userID string, year int, semester int, lang language.Language, course ...string) ([]int, error)
}

type PartialBlueprintAdd = func(hxSwap, hxTarget, hxInclude string, semesters []bool, course string) templ.Component

type Error interface {
	Log(err error)
	Render(w http.ResponseWriter, r *http.Request, code int, userMsg string, lang language.Language)
	RenderPage(w http.ResponseWriter, r *http.Request, code int, userMsg string, title string, userID string, lang language.Language)
}

type Page interface {
	View(main templ.Component, lang language.Language, title string, userID string) templ.Component
}

type Filters interface {
	Init() error
	ParseURLQuery(query url.Values) (expression, error)
	Facets() []string
	IterFacets() any // TODO
}

//================================================================================
// Routing
//================================================================================

func (s *Server) initRouter() {
	router := http.NewServeMux()
	router.HandleFunc("GET /{code}", s.page)
	router.HandleFunc("GET /survey/{code}", s.survey)
	router.HandleFunc("GET /survey/next/{code}", s.surveyNext)
	router.HandleFunc("PUT /rating/{code}/{category}", s.rateCategory)
	router.HandleFunc("DELETE /rating/{code}/{category}", s.deleteCategoryRating)
	router.HandleFunc("PUT /rating/{code}", s.rate)
	router.HandleFunc("DELETE /rating/{code}", s.deleteRating)
	router.HandleFunc("POST /blueprint", s.addCourseToBlueprint)

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
// Handlers
//================================================================================

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID, err := s.Auth.UserID(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.errPageTitle, "", lang)
		return
	}
	code := r.PathValue("code")
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
	s.Page.View(main, lang, title, userID).Render(r.Context(), w)
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
	SurveyFiltersContent(model, texts[model.lang]).Render(r.Context(), w)
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
	SurveysContent(model, texts[model.lang]).Render(r.Context(), w)
}

func (s Server) surveyViewModel(r *http.Request) (surveyViewModel, error) {
	var result surveyViewModel
	req, err := s.parseQueryRequest(r)
	if err != nil {
		return result, errorx.AddContext(err)
	}
	code := r.PathValue("code")
	req.filter.Append("course_code", code)
	searchResponse, err := s.Search.comments(req)
	if err != nil {
		return result, errorx.AddContext(err, errorx.P("code", code))
	}
	result.lang = req.lang
	result.code = code
	result.survey = searchResponse.Survey
	result.offset = req.offset
	result.isEnd = searchResponse.EstimatedTotalHits <= req.offset+req.limit
	result.facets = s.Filters.IterFiltersWithFacets(searchResponse.FacetDistribution, r.URL.Query(), req.lang)
	result.query = req.query
	return result, nil
}

func (s Server) parseQueryRequest(r *http.Request) (request, error) {
	var req request
	userID, err := s.Auth.UserID(r)
	if err != nil {
		return req, errorx.AddContext(err)
	}
	lang := language.FromContext(r.Context())
	query := r.FormValue(searchQuery)
	offset, err := strconv.Atoi(r.FormValue(surveyOffset))
	if err != nil {
		offset = 0
	}
	filter, err := s.Filters.ParseURLQuery(r.URL.Query(), lang)
	if err != nil {
		return req, errorx.AddContext(err)
	}

	req = request{
		userID:   userID,
		query:    query,
		indexUID: "courses-comments", // TODO
		offset:   offset,
		limit:    numberOfComments,
		lang:     lang,
		filter:   expression(&filter),
		facets:   s.Filters.Facets(),
		sort:     "academic_year:desc", // TODO
	}
	return req, nil
}

func (s Server) rateCategory(w http.ResponseWriter, r *http.Request) {
	// get language from context
	lang := language.FromContext(r.Context())
	// get user
	userID, err := s.Auth.UserID(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, texts[lang].errPageTitle, "", lang)
		return
	}
	// get data
	code := r.PathValue("code")
	category := r.PathValue("category")
	ratingString := r.FormValue("rating")
	rating, err := strconv.Atoi(ratingString)
	if err != nil {
		s.Error.Log(errorx.AddContext(err, errorx.P("rating", ratingString)))
		s.Error.Render(w, r, http.StatusBadRequest, texts[lang].errRatingMustBeInt, lang)
		return
	}
	if rating < minRating || rating > maxRating {
		s.Error.Log(errorx.AddContext(fmt.Errorf("rating is not between %d and %d", minRating, maxRating), errorx.P("rating", ratingString)))
		s.Error.Render(w, r, http.StatusBadRequest, texts[lang].errInvalidRating0to10, lang)
		return
	}
	// update category rating
	updatedRating, err := s.Data.rateCategory(userID, code, category, rating, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	// render category rating
	CategoryRating(updatedRating, code, texts[lang]).Render(r.Context(), w)
}

func (s Server) deleteCategoryRating(w http.ResponseWriter, r *http.Request) {
	// get language from context
	lang := language.FromContext(r.Context())
	// get user
	userID, err := s.Auth.UserID(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, texts[lang].errPageTitle, "", lang)
		return
	}
	// get data
	code := r.PathValue("code")
	category := r.PathValue("category")
	// delete category rating
	updatedRating, err := s.Data.deleteCategoryRating(userID, code, category, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	// render category rating
	CategoryRating(updatedRating, code, texts[lang]).Render(r.Context(), w)
}

func (s Server) rate(w http.ResponseWriter, r *http.Request) {
	// get language
	lang := language.FromContext(r.Context())
	// get user
	userID, err := s.Auth.UserID(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, texts[lang].errPageTitle, "", lang)
		return
	}
	// get data
	code := r.PathValue("code")
	ratingString := r.FormValue("rating")
	rating, err := strconv.Atoi(ratingString)
	if err != nil {
		s.Error.Log(errorx.AddContext(err, errorx.P("rating", ratingString)))
		s.Error.Render(w, r, http.StatusBadRequest, texts[lang].errRatingMustBeInt, lang)
		return
	}
	if rating != negativeRating && rating != positiveRating {
		s.Error.Log(errorx.AddContext(fmt.Errorf("rating is not %d or %d", negativeRating, positiveRating), errorx.P("rating", ratingString)))
		s.Error.Render(w, r, http.StatusBadRequest, texts[lang].errInvalidRating0or1, lang)
		return
	}
	// update db
	updatedRating, err := s.Data.rate(userID, code, rating, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	// render the overall rating
	OverallRating(updatedRating, code, texts[lang]).Render(r.Context(), w)
}

func (s Server) deleteRating(w http.ResponseWriter, r *http.Request) {
	// get language
	lang := language.FromContext(r.Context())
	// get user
	userID, err := s.Auth.UserID(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, texts[lang].errPageTitle, "", lang)
		return
	}
	// get code
	code := r.PathValue("code")
	// delete rating
	updatedRating, err := s.Data.deleteRating(userID, code, lang)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.Render(w, r, code, userMsg, lang)
		return
	}
	// render the overall rating
	OverallRating(updatedRating, code, texts[lang]).Render(r.Context(), w)
}

func (s Server) addCourseToBlueprint(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID, err := s.Auth.UserID(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.errPageTitle, "", lang)
		return
	}
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
	Content(&courseDetailPage, t, btn).Render(r.Context(), w)
}

func (s Server) pageNotFound(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID, err := s.Auth.UserID(r)
	if err != nil {
		code, userMsg := errorx.UnwrapError(err, lang)
		s.Error.Log(errorx.AddContext(err))
		s.Error.RenderPage(w, r, code, userMsg, t.errPageTitle, "", lang)
		return
	}
	s.Error.RenderPage(w, r, http.StatusNotFound, t.errPageNotFound, t.errPageTitle, userID, lang)
}
