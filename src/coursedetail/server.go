package coursedetail

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/michalhercik/RecSIS/dbds"
	"github.com/michalhercik/RecSIS/filters"
	"github.com/michalhercik/RecSIS/language"
)

type Filters interface {
	Init() error
	ParseURLQuery(query url.Values) (Expression, error)
	Facets() []string
	IterFacets() any // TODO
}

type Server struct {
	router  *http.ServeMux
	Data    DBManager
	Filters filters.Filters
	Auth    Authentication
	Page    Page
	BpBtn   BlueprintAddButton
	Search  Search
}

//================================================================================
// Interface
//================================================================================

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

//================================================================================
// Init
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
	router.HandleFunc("POST /blueprint/{coursecode}", s.addCourseToBlueprint)
	s.router = router
}

//================================================================================
// Handlers
//================================================================================

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	code := r.PathValue("code")
	course, err := s.course(userID, code, lang)
	if err != nil {
		log.Printf("HandlePage error %s: %v", code, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	numberOfYears, err := s.BpBtn.NumberOfYears(userID)
	if err != nil {
		http.Error(w, "Unable to retrieve number of years", http.StatusInternalServerError)
		log.Printf("HandlePage error: %v", err)
		return
	}
	btn := s.BpBtn.PartialComponent(numberOfYears, lang)
	main := Content(course, t, btn)
	s.Page.View(main, lang, course.Code+" - "+course.Name).Render(r.Context(), w)
}

func (s Server) course(userID, code string, lang language.Language) (*Course, error) {
	var result *Course
	result, err := s.Data.Course(userID, code, lang)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s Server) surveyNext(w http.ResponseWriter, r *http.Request) {
	model, err := s.surveyViewModel(r)
	if err != nil {
		log.Printf("survey %s: %v", model.code, err)
		http.Error(w, "Unable to parse request", http.StatusBadRequest)
		return
	}
	SurveysContent(model, texts[model.lang]).Render(r.Context(), w)
}

func (s Server) survey(w http.ResponseWriter, r *http.Request) {
	model, err := s.surveyViewModel(r)
	if err != nil {
		log.Printf("survey %s: %v", model.code, err)
		http.Error(w, "Unable to parse request", http.StatusBadRequest)
		return
	}
	Survey(model, texts[model.lang]).Render(r.Context(), w)
}

func (s Server) surveyViewModel(r *http.Request) (SurveyViewModel, error) {
	var result SurveyViewModel
	code := r.PathValue("code")
	req, err := s.parseQueryRequest(r)
	if err != nil {
		return result, err
	}
	req.filter.Append("course_code", code)
	searchResponse, err := s.Search.Comments(req)
	if err != nil {
		return result, err
	}
	result.lang = req.lang
	result.code = code
	result.survey = searchResponse.Survey
	result.isEnd = searchResponse.EstimatedTotalHits <= req.offset+req.limit
	result.facets = s.Filters.IterFiltersWithFacets(searchResponse.FacetDistribution, r.URL.Query(), req.lang)
	result.query = req.query
	return result, nil
}

func (s Server) rate(w http.ResponseWriter, r *http.Request) {
	// get the course code from the request
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	code := r.PathValue("code")
	rating, err := strconv.Atoi(r.FormValue("rating"))
	if err != nil && (rating == positiveRating || rating == negativeRating) {
		log.Printf("rate error: %v", err)
		return
	}
	updatedRating, err := s.Data.Rate(userID, code, rating)
	if err != nil {
		log.Printf("rate error: %v", err)
	}
	// get language
	lang := language.FromContext(r.Context())
	// render the overall rating
	OverallRating(updatedRating.UserRating, updatedRating.AvgRating, updatedRating.RatingCount, code, texts[lang]).Render(r.Context(), w)
}

func (s Server) deleteRating(w http.ResponseWriter, r *http.Request) {
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	code := r.PathValue("code")
	updatedRating, err := s.Data.DeleteRating(userID, code)
	if err != nil {
		log.Printf("deleteRating error: %v", err)
	}
	// get language
	lang := language.FromContext(r.Context())
	// render the overall rating
	OverallRating(updatedRating.UserRating, updatedRating.AvgRating, updatedRating.RatingCount, code, texts[lang]).Render(r.Context(), w)
}

func (s Server) rateCategory(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	code := r.PathValue("code")
	category := r.PathValue("category")
	rating, err := strconv.Atoi(r.FormValue("rating"))
	if err != nil {
		log.Printf("rateCategory error: %v", err)
		return
	}
	//TODO handle language properly
	updatedRating, err := s.Data.RateCategory(userID, code, category, rating, lang)
	if err != nil {
		log.Printf("rateCategory error: %v", err)
	}
	// render category rating
	CategoryRating(updatedRating, code, texts[lang]).Render(r.Context(), w)
}

func (s Server) deleteCategoryRating(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	code := r.PathValue("code")
	category := r.PathValue("category")
	//TODO handle language properly
	updatedRating, err := s.Data.DeleteCategoryRating(userID, code, category, lang)
	if err != nil {
		log.Printf("deleteCategoryRating error: %v", err)
	}
	// render category rating
	CategoryRating(updatedRating, code, texts[lang]).Render(r.Context(), w)
}

func (s Server) addCourseToBlueprint(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	courseCode := r.PathValue("coursecode")
	year, err := strconv.Atoi(r.FormValue("year"))
	if err != nil {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	semesterInt, err := strconv.Atoi(r.FormValue("semester"))
	if err != nil {
		http.Error(w, "Invalid semester", http.StatusBadRequest)
		return
	}
	semester := dbds.SemesterAssignment(semesterInt)
	_, err = s.BpBtn.Action(userID, year, semester, courseCode)
	if err != nil {
		http.Error(w, "Unable to add course to blueprint", http.StatusInternalServerError)
		log.Printf("HandlePage error: %v", err)
		return
	}
	numberOfYears, err := s.BpBtn.NumberOfYears(userID)
	if err != nil {
		http.Error(w, "Unable to retrieve number of years", http.StatusInternalServerError)
		log.Printf("HandlePage error: %v", err)
		return
	}
	t := texts[lang]
	btn := s.BpBtn.PartialComponent(numberOfYears, lang)
	course, err := s.course(userID, courseCode, lang)
	if err != nil {
		log.Printf("HandlePage error %s: %v", courseCode, err)
		http.Error(w, "Course not found", http.StatusNotFound)
	} else {
		main := Content(course, t, btn)
		s.Page.View(main, lang, course.Code+" - "+course.Name).Render(r.Context(), w)
	}
}

func (s Server) parseQueryRequest(r *http.Request) (Request, error) {
	var req Request
	userID, err := s.Auth.UserID(r)
	if err != nil {
		return req, err
	}
	lang := language.FromContext(r.Context())
	query := r.FormValue("survey-search")
	offset, err := strconv.Atoi(r.FormValue(nOCommentsQuery))
	if err != nil {
		offset = 1
	}
	filter, err := s.Filters.ParseURLQuery(r.URL.Query())
	if err != nil {
		// TODO: handle error
		log.Printf("search error: %v", err)
	}

	req = Request{
		userID:   userID,
		query:    query,
		indexUID: "courses-comments", // TODO
		offset:   offset,
		limit:    numberOfComments,
		lang:     lang,
		filter:   Expression(&filter),
		facets:   s.Filters.Facets(),
		sort:     "academic_year:desc", // TODO
	}
	return req, nil
}
