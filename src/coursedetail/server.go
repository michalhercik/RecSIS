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
	Auth           Authentication
	Page           Page
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

func (s Server) page(w http.ResponseWriter, r *http.Request) {
	lang := language.FromContext(r.Context())
	t := texts[lang]
	userID, err := s.Auth.UserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	code := r.PathValue("code")
	course, err := s.course(userID, code, lang, r)
	if err != nil {
		log.Printf("HandlePage error %s: %v", code, err)
		PageNotFound(code, t).Render(r.Context(), w)
	} else {
		main := Content(course, t)
		s.Page.View(main, lang, course.Code+" - "+course.Name).Render(r.Context(), w)
	}
}

func (s Server) course(userID, code string, lang language.Language, r *http.Request) (*Course, error) {
	var result *Course
	result, err := s.Data.Course(userID, code, lang)
	if err != nil {
		return nil, err
	}
	br := s.CourseComments.BuildRequest(lang)
	br, err = br.ParseURLQuery(r.URL.Query())
	if err != nil {
		return nil, err
	}
	// TODO: jako GitHub searchbar, filters
	// TODO: paginace prev/next (nic vic)
	// TODO: alterative -> only load more (hx-trigger revealed) lazy loading
	searchQuery := r.FormValue("q") // TODO: input text name=q (RENAME)
	br = br.SetQuery(searchQuery)
	br = br.AddCourse(code)
	br = br.SetLimit(20) // TODO: pocet komentaru na stranku
	br = br.SetOffset(0) // TODO: offset komentaru
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
