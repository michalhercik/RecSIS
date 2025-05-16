package coursedetail

import (
	"log"
	"net/http"
	"strconv"

	"github.com/michalhercik/RecSIS/dbds"
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
	BpBtn          BlueprintAddButton
}

func (s *Server) Init() {
	router := http.NewServeMux()
	router.HandleFunc("GET /{code}", s.page)
	router.HandleFunc("GET /survey/{code}", s.survey)
	router.HandleFunc("PUT /rating/{code}/{category}", s.rateCategory)
	router.HandleFunc("DELETE /rating/{code}/{category}", s.deleteCategoryRating)
	router.HandleFunc("PUT /rating/{code}", s.rate)
	router.HandleFunc("DELETE /rating/{code}", s.deleteRating)
	router.HandleFunc("POST /blueprint/{coursecode}", s.addCourseToBlueprint)
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
	br = br.SetQuery("") // default query is empty
	br = br.AddCourse(code)
	br = br.SetLimit(numberOfComments)
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

func (s Server) survey(w http.ResponseWriter, r *http.Request) {
	// language
	lang := language.FromContext(r.Context())
	// code
	code := r.PathValue("code")
	// number of comments
	numberOfCommentsStr := r.FormValue(nOCommentsQuery)
	numberOfComments, err := strconv.Atoi(numberOfCommentsStr)
	if err != nil {
		log.Printf("HandlePage error %s: %v", code, err)
	}
	// survey
	br := s.CourseComments.BuildRequest(lang)
	br, err = br.ParseURLQuery(r.URL.Query())
	if err != nil {
		log.Printf("HandlePage error %s: %v", code, err)
	}
	// TODO: paginace prev/next (nic vic)
	// TODO: alterative -> only load more (hx-trigger revealed) lazy loading
	// TODO: if alterative -> add top button
	searchQuery := r.FormValue("survey-search")
	br = br.SetQuery(searchQuery)
	br = br.AddCourse(code)
	br = br.SetLimit(numberOfComments)
	br = br.SetOffset(0)
	br = br.AddSort(params.AcademicYear, params.Desc)
	req, err := br.Build()
	if err != nil {
		log.Printf("HandlePage error %s: %v", code, err)
	}
	surveyRes, err := s.CourseComments.Comments(req)
	if err != nil {
		log.Printf("HandlePage error %s: %v", code, err)
	}
	survey := surveyRes.Comments()
	noMoreComments := len(survey) < numberOfComments
	// render the survey
	SurveysContent(survey, code, noMoreComments, texts[lang]).Render(r.Context(), w)
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
	_, err = s.BpBtn.Action(userID, courseCode, year, semester)
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
	// course, err := s.Data.Course(userID, courseCode, lang)
	// if err != nil {
	// 	http.Error(w, "Unable to retrieve course", http.StatusInternalServerError)
	// 	log.Printf("HandlePage error: %v", err)
	// 	return
	// }
	// CourseRow(&course, btn, t).Render(r.Context(), w)
	course, err := s.course(userID, courseCode, lang, r)
	if err != nil {
		log.Printf("HandlePage error %s: %v", courseCode, err)
		http.Error(w, "Course not found", http.StatusNotFound)
	} else {
		main := Content(course, t, btn)
		s.Page.View(main, lang, course.Code+" - "+course.Name).Render(r.Context(), w)
	}
}
