package coursedetail

import (
	"log"
	"net/http"
	"strconv"

	"github.com/michalhercik/RecSIS/internal/course/comments/meilisearch/params"
	"github.com/michalhercik/RecSIS/internal/course/comments/search"
	"github.com/michalhercik/RecSIS/language"
)

type Server struct {
	router         *http.ServeMux
	Data           DataManager
	CourseComments search.SearchEngine
}

func (s *Server) Init() {
	router := http.NewServeMux()
	router.HandleFunc("GET /{code}", s.page)
	router.HandleFunc("PUT /rating/{code}/{category}", s.rateCategory)
	router.HandleFunc("DELETE /rating/{code}/{category}", s.deleteCategoryRating)
	router.HandleFunc("PUT /rating/{code}", s.rate)
	router.HandleFunc("DELETE /rating/{code}", s.deleteRating)
	s.router = router
}

func (s Server) Router() http.Handler {
	return s.router
}

// func (s Server) Register(router *http.ServeMux, prefix string) {
// 	lr := language.LanguageRouter{Router: router}
// 	lr.HandleLangFunc(fmt.Sprintf("%s/{code}", prefix), http.MethodGet, s.page)
// 	lr.HandleLangFunc(fmt.Sprintf("%s/rating/{code}/{category}", prefix), http.MethodPut, s.rateCategory)
// 	lr.HandleLangFunc(fmt.Sprintf("%s/rating/{code}/{category}", prefix), http.MethodDelete, s.deleteCategoryRating)
// 	lr.HandleLangFunc(fmt.Sprintf("%s/rating/{code}", prefix), http.MethodPut, s.rate)
// 	lr.HandleLangFunc(fmt.Sprintf("%s/rating/{code}", prefix), http.MethodDelete, s.deleteRating)
// }

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	log.Println("Jsem tady!")
	log.Println(r.URL.Path)
	lang := language.FromContext(r.Context())
	t := texts[lang]
	sessionCookie, err := r.Cookie("recsis_session_key")
	if err != nil {
		log.Printf("courseDetail error: %v", err)
		return
	}
	code := r.PathValue("code")
	course, err := s.course(sessionCookie.Value, code, lang, r)
	if err != nil {
		log.Printf("HandlePage error %s: %v", code, err)
		PageNotFound(code, t).Render(r.Context(), w)
	} else {
		Page(course, t).Render(r.Context(), w)
	}
}

func (s Server) course(sessionID, code string, lang language.Language, r *http.Request) (*Course, error) {
	var result *Course
	result, err := s.Data.Course(sessionID, code, lang)
	if err != nil {
		return nil, err
	}
	br := s.CourseComments.BuildRequest(lang)
	br, err = br.ParseURLQuery(r.URL.Query())
	if err != nil {
		return nil, err
	}
	searchQuery := r.FormValue("q")
	br = br.SetQuery(searchQuery)
	br = br.AddCourse(code)
	br = br.SetLimit(20)
	br = br.SetOffset(0)
	br = br.AddSort(params.AcademicYear, params.Desc)
	req, err := br.Build()
	if err != nil {
		return nil, err
	}
	result.Comments, err = s.CourseComments.Comments(req)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s Server) rate(w http.ResponseWriter, r *http.Request) {
	// get the course code from the request
	sessionCookie, err := r.Cookie("recsis_session_key")
	if err != nil {
		log.Printf("rate error: %v", err)
		return
	}
	code := r.PathValue("code")
	rating, err := strconv.Atoi(r.FormValue("rating"))
	if err != nil && (rating == positiveRating || rating == negativeRating) {
		log.Printf("rate error: %v", err)
		return
	}
	updatedRating, err := s.Data.Rate(sessionCookie.Value, code, rating)
	if err != nil {
		log.Printf("rate error: %v", err)
	}
	// TODO: render overall rating
	_ = updatedRating
}

func (s Server) deleteRating(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("recsis_session_key")
	if err != nil {
		log.Printf("deleteRating error: %v", err)
		return
	}
	code := r.PathValue("code")
	updatedRating, err := s.Data.DeleteRating(sessionCookie.Value, code)
	if err != nil {
		log.Printf("deleteRating error: %v", err)
	}
	// TODO: render overall rating
	_ = updatedRating
}

func (s Server) rateCategory(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	sessionCookie, err := r.Cookie("recsis_session_key")
	if err != nil {
		log.Printf("rateCategory error: %v", err)
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
	updatedRating, err := s.Data.RateCategory(sessionCookie.Value, code, category, rating, lang)
	if err != nil {
		log.Printf("rateCategory error: %v", err)
	}
	// TOOD: render category rating
	_ = updatedRating
}

func (s Server) deleteCategoryRating(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	sessionCookie, err := r.Cookie("recsis_session_key")
	if err != nil {
		log.Printf("deleteCategoryRating error: %v", err)
		return
	}
	code := r.PathValue("code")
	category := r.PathValue("category")
	//TODO handle language properly
	updatedRating, err := s.Data.DeleteCategoryRating(sessionCookie.Value, code, category, lang)
	if err != nil {
		log.Printf("deleteCategoryRating error: %v", err)
	}
	// TOOD: render category rating
	_ = updatedRating
}
